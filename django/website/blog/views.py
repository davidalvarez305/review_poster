from datetime import date
import math
import os
from django.shortcuts import render, get_object_or_404
from django.views import View
from django.db import connection
from .models import *
from django.contrib.auth.mixins import LoginRequiredMixin
from django.http import HttpResponseBadRequest, JsonResponse
from django.contrib.auth import authenticate, login, logout

class MyBaseView(View):
    groups = Group.prefetch_category_set()
    domain = str(os.environ.get('DJANGO_DOMAIN'))
    current_year = date.today().year
    google_analytics_id = str(os.environ.get('GOOGLE_ANALYTICS_ID'))

    context = {
        'domain': domain,
        'current_year': current_year,
        'groups': groups,
        'google_analytics_id': google_analytics_id,
        'google_analytics_src': "https://www.googletagmanager.com/gtag/js?id=" + google_analytics_id,
        'meta_description': 'Get reviews for all things sports, fitness, outdoors, and everything in between!',
        'page_title': str(os.environ.get('SITE_NAME')),
        'site_name': str(os.environ.get('SITE_NAME')),
        'is_reviewpost': False,
    }

    template_name = 'home.html'

    def get(self, request, *args, **kwargs):
        return render(request, self.template_name, context=self.context)


class HomeView(MyBaseView):
    template_name = 'blog/home.html'

    def get(self, request, *args, **kwargs):
        context = self.context
        context['page_path'] = request.build_absolute_uri()
        context['page_title'] = str(os.environ.get('SITE_NAME'))
        return render(request, self.template_name, context=context)


class CategoryView(MyBaseView):
    template_name = 'blog/category.html'
        
    def get(self, request, *args, **kwargs):
        context = self.context
        category_slug = kwargs['category']
        category = get_object_or_404(Category, slug=category_slug)
        sub_categories = SubCategory.objects.filter(category__slug=category_slug).prefetch_related('reviewpost_set')
        context['sub_categories'] = sub_categories
        context['category'] = category
        context['page_title'] = category.name.title() + " - " + str(os.environ.get('SITE_NAME'))
        context['page_path'] = request.build_absolute_uri()

        return render(request, self.template_name, context)


class SubCategoryView(MyBaseView):
    template_name = 'blog/sub_category.html'

    def get(self, request, *args, **kwargs):
        context = self.context
        sub_category_slug = kwargs['sub_category']
        sub_category = get_object_or_404(SubCategory, slug=sub_category_slug)
        posts = ReviewPost.objects.filter(sub_category__slug=sub_category_slug)
        context['category_slug'] = kwargs['category']
        context['posts'] = posts
        context['sub_category'] = sub_category
        context['page_title'] = sub_category.name.title() + " - " + str(os.environ.get('SITE_NAME'))
        context['page_path'] = request.build_absolute_uri()
        return render(request, self.template_name, context)

class ReviewPostView(MyBaseView):
    template_name = 'blog/review_post.html'

    def get(self, request, *args, **kwargs):
        context = self.context
        slug = kwargs['slug']
        sub_category_slug = kwargs['sub_category']

        review_post = get_object_or_404(ReviewPost.objects.select_related('sub_category'), slug=slug)
        related_review_posts = ReviewPost.objects.filter(sub_category__slug=sub_category_slug)

        i = 0
        product_rating_stars = ""
        if review_post.product.product_ratings == "":
            product_rating_stars = 0
        else:
            while i < float(review_post.product.product_ratings):
                product_rating_stars += '<svg class="hi-mini hi-star inline-block w-5 h-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true"><path fill-rule="evenodd" d="M10.868 2.884c-.321-.772-1.415-.772-1.736 0l-1.83 4.401-4.753.381c-.833.067-1.171 1.107-.536 1.651l3.62 3.102-1.106 4.637c-.194.813.691 1.456 1.405 1.02L10 15.591l4.069 2.485c.713.436 1.598-.207 1.404-1.02l-1.106-4.637 3.62-3.102c.635-.544.297-1.584-.536-1.65l-4.752-.382-1.831-4.401z" clip-rule="evenodd"/></svg>'
                i += 1
        if len(review_post.product.product_price) == 0:
            review_post.product.product_price = '$0.00'
        context['related_review_posts'] = related_review_posts
        context['review_post'] = review_post
        context['meta_description'] = review_post.description
        context['page_title'] = review_post.title + " - " + str(os.environ.get('SITE_NAME'))
        context['product_rating_stars'] = product_rating_stars
        context['sub_category_slug'] = kwargs['sub_category']
        context['category_slug'] = kwargs['category']
        context['page_path'] = request.build_absolute_uri()
        context['is_reviewpost'] = True
        return render(request, self.template_name, context)

def sitemap(request, *args, **kwargs):

    sitemap_index = kwargs['int'] - 1

    if sitemap_index < 0:
        return HttpResponseBadRequest("Sitemap value must be greater than zero.")

    offset = 5000 * sitemap_index

    cursor = connection.cursor()
    cursor.execute(
        f'''SELECT CONCAT('0.8'), review_post.slug AS slug, sc.slug AS sub_category_slug, c.slug AS category_slug
            FROM review_post
            LEFT JOIN sub_category AS sc
            ON review_post.sub_category_id = sc.id
            LEFT JOIN category AS c
            ON sc.category_id = c.id
            LIMIT 5000 OFFSET {offset};''')
    rows = cursor.fetchall()
    columns = ["priority", "slug", "sub_category_slug", "category_slug"]
    posts = [
        dict(zip(columns, row))
        for row in rows
    ]

    if len(posts) == 0:
        return HttpResponseBadRequest("Bad Request: Exceeded post quantity.")

    context = {
        'posts': posts,
        'domain': str(os.environ.get('DJANGO_DOMAIN')),
        "current_year": date.today().year
    }
    return render(request, 'blog/sitemap.xml.gz', context, content_type="application/xhtml+xml")


def sitemap_index(request, *args, **kwargs):

    cursor = connection.cursor()
    cursor.execute(
        '''SELECT CONCAT('0.8'), review_post.slug AS slug FROM review_post;
        ''')
    rows = cursor.fetchall()
    columns = ["priority", "slug"]
    posts = [
        dict(zip(columns, row))
        for row in rows
    ]

    # Helpers for slicing sitemaps into chunks of 5000 (Google MAX: 50,0000)
    start = 0
    interval = 5000
    i = 0

    # Amount of times to slice the list
    loops = math.ceil(len(posts) / interval)
    sitemaps = []

    while (i < loops):
        sliced_sitemap = posts[start:start+interval]
        sitemaps.append(sliced_sitemap)
        start += 5000
        i += 1

    context = {
        'sitemaps': sitemaps,
        'domain': str(os.environ.get('DJANGO_DOMAIN')),
        "current_year": date.today().year
    }
    return render(request, 'blog/sitemap_index.xml', context, content_type="application/xhtml+xml")

class AffiliateDisclaimer(MyBaseView):
    template_name = 'blog/affiliate_disclaimer.html'

    def get(self, request, *args, **kwargs):
        context = self.context
        context['page_title'] = "Affiliate Disclaimer - " + context['site_name']
        return render(request, self.template_name, context)

class PrivacyPolicy(MyBaseView):
    template_name = 'blog/privacy_policy.html'

    def get(self, request, *args, **kwargs):
        context = self.context
        context['page_title'] = "Privacy Policy - " + context['site_name']
        return render(request, self.template_name, context)