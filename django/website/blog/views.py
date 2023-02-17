from datetime import date
import os
from django.shortcuts import render, get_object_or_404
from django.views import View
from .models import *

class MyBaseView(View):
    groups = Group.objects.all()
    domain = str(os.environ.get('DOMAIN'))
    current_year = date.today().year

    context = {
        'domain': domain,
        'current_year': current_year,
        'groups': groups,
    }

    template_name = 'home.html'

    def get(self, request, *args, **kwargs):
        return render(request, self.template_name, context=self.context)


class HomeView(MyBaseView):
    template_name = 'blog/home.html'


class CategoryView(MyBaseView):

    template_name = 'blog/category.html'

    def get_context_data(self, **kwargs):
        context = super().get_context_data(**kwargs)
        category_slug = kwargs['category']
        category = get_object_or_404(Category, slug=category_slug)
        related_sub_categories = SubCategory.objects.filter(category__slug=category_slug)
        context['related_sub_categories'] = related_sub_categories
        return context


class SubCategoryView(MyBaseView):

    template_name = 'blog/sub_category.html'

    def get_context_data(self, **kwargs):
        context = super().get_context_data(**kwargs)
        sub_category_slug = kwargs['sub_category']
        category_slug = kwargs['category']
        sub_category = get_object_or_404(SubCategory, sub_category_slug=sub_category_slug)
        related_sub_categories = SubCategory.objects.filter(category__slug=category_slug)
        context['related_sub_categories'] = related_sub_categories
        context['sub_category'] = sub_category
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
        return context