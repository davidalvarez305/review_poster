from datetime import date
import os
from django.shortcuts import render, get_object_or_404
from .models import *

# Environment Variables
DOMAIN = str(os.environ.get('DOMAIN'))
CURRENT_YEAR = date.today().year

def home(request):
    context = {
        "domain": DOMAIN,
        "current_year": CURRENT_YEAR
    }
    return render(request, 'blog/home.html', context)

def review_post(request, **kwargs):
    sub_category_slug = kwargs['sub_category']
    slug = kwargs['slug']
    post = get_object_or_404(ReviewPost, slug=slug)
    sub_category = get_object_or_404(SubCategory, slug=sub_category_slug)
    related_posts = ReviewPost.objects.filter(sub_category__slug=sub_category_slug)

    context = {
        "domain": DOMAIN,
        "current_year": CURRENT_YEAR,
        "post": post,
        "related_posts": related_posts,
        "sub_category": sub_category
    }

    return render(request, 'blog/review_post.html', context)