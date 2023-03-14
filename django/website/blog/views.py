from datetime import date
import math
import os
from django.http import HttpResponseBadRequest
from django.shortcuts import render, get_object_or_404
from django.views import View
from django.db import connection
from .models import *
from django.db.models import Prefetch

class MyBaseView(View):
    groups = Group.objects.prefetch_related('category_set').all()
    domain = str(os.environ.get('DOMAIN'))
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
        'site_name': str(os.environ.get('SITE_NAME'))
    }

    template_name = 'home.html'

    def get(self, request, *args, **kwargs):
        ctx = self.context
        ctx['path'] = request.path
        return render(request, self.template_name, context=ctx)


class HomeView(MyBaseView):
    template_name = 'blog/home.html'

    def get_context_data(self, **kwargs):
        cursor = connection.cursor()
        context = super().get_context_data(**kwargs)

        sql = '''
            SELECT JSON_AGG(TO_JSON(r)) FROM (
                SELECT rp.title, rp.slug, rp.horizontalcardproductimageurl, rp."categoryId", sub.slug AS sub_category_slug, sub.name AS sub_category, cg.name AS category, cg.id AS category_id,
                ROW_NUMBER() OVER (PARTITION BY cg.id ORDER BY cg.id) row_number
                FROM review_post AS rp
                LEFT JOIN sub_category AS sub
                ON rp."categoryId" = sub.id
                LEFT JOIN category AS cg
                ON cg.id = sub.category_id
                GROUP BY sub.id, rp.title, sub.name, sub.slug, cg.name, rp.slug, rp.horizontalcardproductimageurl, cg.id, rp."categoryId"
                ORDER BY cg.id
            ) AS r
            WHERE r.row_number < 10
        '''

        cursor.execute(sql)
        rows = cursor.fetchall()
        desc = cursor.description
        columns = [col[0] for col in desc]
        example_posts = [
            dict(zip(columns, row))
            for row in rows
        ]

        context['example_posts'] = example_posts
        return context


class CategoryView(MyBaseView):
    template_name = 'blog/category.html'
        
    def get(self, request, *args, **kwargs):
        context = self.context
        category_slug = kwargs['category']
        category = get_object_or_404(Category, slug=category_slug)
        sub_categories = SubCategory.objects.filter(category__slug=category_slug)
        example_posts = sub_categories.prefetch_related(Prefetch('review_post', queryset=ReviewPost.objects.first()))
        context['example_posts'] = example_posts
        context['sub_categories'] = sub_categories
        context['category'] = category
        context['page_title'] = category.name.title() + " - " + str(os.environ.get('SITE_NAME'))
        return render(request, self.template_name, context)


class SubCategoryView(MyBaseView):
    template_name = 'blog/sub_category.html'

    def get(self, request, *args, **kwargs):
        context = self.context
        sub_category_slug = kwargs['sub_category']
        sub_category = get_object_or_404(SubCategory, sub_category_slug=sub_category_slug)
        posts = ReviewPost.objects.filter(sub_category__slug=sub_category.slug)
        context['posts'] = posts
        context['sub_category'] = sub_category
        context['page_title'] = sub_category.name.title() + " - " + str(os.environ.get('SITE_NAME'))
        return render(request, self.template_name, context)

class ReviewPostView(MyBaseView):
    template_name = 'blog/review_post.html'

    def get(self, request, *args, **kwargs):
        context = self.context
        slug = kwargs['slug']
        sub_category_slug = kwargs['sub_category']

        review_post = get_object_or_404(ReviewPost, slug=slug)
        product = get_object_or_404(Product, affiliate_url=review_post.product_affiliate_url)

        i = 0
        product_rating_stars = ""
        if product.product_ratings == "":
            product_rating_stars = 0
        else:
            while i < float(product.product_ratings):
                product_rating_stars += '<li><i class="flaticon-star"></i></li>'
                i += 1

        related_review_posts = ReviewPost.objects.filter(sub_category__slug=sub_category_slug)
        context['related_review_posts'] = related_review_posts
        context['review_post'] = review_post
        context['meta_description'] = review_post.description
        context['page_title'] = review_post.title
        context['product'] = product
        context['product_rating_stars'] = product_rating_stars
        return render(request, self.template_name, context)

def sitemap(request, *args, **kwargs):

    sitemap_index = kwargs['int'] - 1

    if sitemap_index < 0:
        return HttpResponseBadRequest("Sitemap value must be greater than zero.")

    offset = 5000 * sitemap_index

    cursor = connection.cursor()
    cursor.execute(
        f"SELECT CONCAT('0.8'), slug FROM review_post LIMIT 5000 OFFSET {offset};")
    rows = cursor.fetchall()
    columns = ["priority", "slug"]
    posts = [
        dict(zip(columns, row))
        for row in rows
    ]

    if len(posts) == 0:
        return HttpResponseBadRequest("Bad Request: Exceeded post quantity.")

    context = {
        'posts': posts,
        'domain': str(os.environ.get('DOMAIN')),
        "current_year": date.today().year
    }
    return render(request, 'blog/sitemap.xml.gz', context, content_type="application/xhtml+xml")


def sitemap_index(request, *args, **kwargs):

    cursor = connection.cursor()
    cursor.execute(
        '''SELECT CONCAT('0.8'), slug, sc.slug AS sub_category_slug FROM review_post
            LEFT JOIN sub_category AS sc
            ON review_post.sub_category_id = sc.id;
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
        'domain': str(os.environ.get('DOMAIN')),
        "current_year": date.today().year
    }
    return render(request, 'blog/sitemap_index.xml', context, content_type="application/xhtml+xml")

class AffiliateDisclaimer(MyBaseView):
    template_name = 'blog/affiliate_disclaimer.html'

class PrivacyPolicy(MyBaseView):
    template_name = 'blog/privacy_policy.html'

class CreatePost(MyBaseView):
    template_name = 'blog/create_post.html'

    def get(self, request, *args, **kwargs):
        context = self.context

        options = self.context['groups']
        select_options = ""
        for option in options:
            select_options += f'<option value={option}>{option}</option>'

        context['crawler_api'] = os.environ.get('REVIEW_POST_API') + "/api/review-post"
        context['select_options'] = select_options
        context['page_title'] = "Create Review Posts - " + context['site_name']
        return render(request, self.template_name, context)