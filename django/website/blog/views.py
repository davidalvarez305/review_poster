from datetime import date
import os
from django.shortcuts import render, get_object_or_404
from django.views import View
from django.db import connection
from .models import *

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
        'page_title': str(os.environ.get('PAGE_TITLE')),
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

    def get_context_data(self, **kwargs):
        context = super().get_context_data(**kwargs)
        category_slug = kwargs['category']
        category = get_object_or_404(Category, slug=category_slug)
        sub_categories = SubCategory.objects.filter(category__slug=category_slug)
        context['sub_categories'] = sub_categories
        context['category'] = category
        context['page_title'] = category.name.title() + " - " + str(os.environ.get('PAGE_TITLE'))
        return context


class SubCategoryView(MyBaseView):

    template_name = 'blog/sub_category.html'

    def get_context_data(self, **kwargs):
        context = super().get_context_data(**kwargs)
        sub_category_slug = kwargs['sub_category']
        sub_category = get_object_or_404(SubCategory, sub_category_slug=sub_category_slug)
        posts = ReviewPost.objects.filter(sub_category__slug=sub_category.slug)
        context['posts'] = posts
        context['sub_category'] = sub_category
        context['page_title'] = sub_category.name.title() + " - " + str(os.environ.get('PAGE_TITLE'))
        return context

class ReviewPostView(MyBaseView):

    template_name = 'blog/review_post.html'

    def get_context_data(self, **kwargs):
        context = super().get_context_data(**kwargs)
        slug = kwargs['slug']
        sub_category_slug = kwargs['sub_category']
        review_post = get_object_or_404(ReviewPost, slug=slug)
        related_review_posts = ReviewPost.objects.filter(sub_category__slug=sub_category_slug)
        context['related_review_posts'] = related_review_posts
        context['review_post'] = review_post
        context['meta_description'] = review_post.description
        context['page_title'] = review_post.title
        return context