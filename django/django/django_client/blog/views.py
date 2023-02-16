from datetime import date
import math
from django.http import HttpResponseBadRequest, HttpResponseNotFound
from django.shortcuts import render
from django.db import connection
import os
from .utils import define_silo, loop_items

# Environment Variables
domain = str(os.environ.get('DOMAIN'))
site_name = str(os.environ.get('SITE_NAME'))
logo = " ".join(site_name.lower())
current_year = date.today().year


def home(request):
    cursor = connection.cursor()

    cursor.execute(
        f'''
        SELECT JSON_BUILD_OBJECT('example_posts', (
            SELECT JSON_AGG(TO_JSON(r)) FROM (
                SELECT rp.title, rp.slug, rp.horizontalcardproductimageurl, rp."categoryId", c.name AS category, cg.name AS category_group, cg.id AS category_group_id,
                ROW_NUMBER() OVER (PARTITION BY cg.id ORDER BY cg.id) row_number
                FROM review_post AS rp
                LEFT JOIN category AS c
                ON rp."categoryId" = c.id
                LEFT JOIN category_group AS cg
                ON cg.id = c.category_group_id
                GROUP BY c.id, rp.title, c.name, cg.name, rp.slug, rp.horizontalcardproductimageurl, cg.id, rp."categoryId"
                ORDER BY cg.id
            ) AS r
            WHERE r.row_number < 10
        ),'category', (
            SELECT JSON_AGG(TO_JSON(category_group)) FROM category_group
        ), 'parent_group', (
            SELECT JSON_AGG(TO_JSON(parent_group)) FROM parent_group
        ), 'category_group', (
            SELECT JSON_AGG(TO_JSON(category_group)) FROM category_group
        )) AS home;
        '''
    )
    rows = cursor.fetchall()
    desc = cursor.description
    columns = [col[0] for col in desc]
    cats = [
        dict(zip(columns, row))
        for row in rows
    ]
    meta_description = f"Welcome to the {site_name}."

    example_posts = cats[0]['home']['example_posts']
    categories = cats[0]['home']['category']
    separated_categories = loop_items(items=categories)

    ctgs = []

    for sc in separated_categories:
        sub_cats = []
        for cg in sc:
            for post in example_posts:
                if post['category_group_id'] == cg['id']:
                    cg['horizontalcardproductimageurl'] = post['horizontalcardproductimageurl']
                    sub_cats.append(cg)
                    break
        ctgs.append(sub_cats)

    parent_group = define_silo(cats, 'home')

    google_analytics_id = os.environ.get("GOOGLE_ANALYTICS_ID")

    context = {
        "categories": ctgs,
        "parent_group": parent_group,
        "current_year": current_year,
        "domain": domain,
        "site_name": site_name,
        "logo": logo,
        "page_title": site_name,
        "meta_description": meta_description,
        "example_posts": example_posts[:20],
        "google_analytics_id": google_analytics_id,
        "google_analytics_src": "https://www.googletagmanager.com/gtag/js?id=" + google_analytics_id
    }

    return render(request, 'blog/home.html', context)


def category_group(request, *args, **kwargs):
    category_group_slug = kwargs['name']

    cursor = connection.cursor()
    cursor.execute(
        f'''
        SELECT JSON_BUILD_OBJECT('category', (
            SELECT JSON_AGG(TO_JSON(category)) AS c FROM category WHERE category_group_id = (
                SELECT id FROM category_group AS cg WHERE cg.name = '{category_group_slug}'
            )
        ), 'parent_group', (
            SELECT JSON_AGG(TO_JSON(parent_group)) FROM parent_group
        ), 'category_group', (
            SELECT JSON_AGG(TO_JSON(category_group)) FROM category_group
        ), 'example_post', (
            SELECT JSON_AGG(TO_JSON(rp)) FROM review_post AS rp WHERE rp."categoryId" IN (
                SELECT id FROM category AS c WHERE c.category_group_id = (
                    SELECT id FROM category_group AS cg WHERE cg.name = '{category_group_slug}'
                )
            ) LIMIT 1
        )) AS category;
        '''
    )
    rows = cursor.fetchall()
    desc = cursor.description
    columns = [col[0] for col in desc]
    cats = [
        dict(zip(columns, row))
        for row in rows
    ]

    if not cats[0]['category']['category']:
        return HttpResponseNotFound("Not found.")

    path = request.path.split('/category-group/')[1].replace("-", " ")
    page_title = "Reviews of " + path.title() + " Products | " + site_name
    meta_description = f"Learn more about {path} products & choices with the {site_name}."

    example_posts = cats[0]['category']['example_post']

    categories = cats[0]['category']['category']

    separated_categories = loop_items(items=categories)

    ctgs = []

    for sc in separated_categories:
        sub_cats = []
        for cg in sc:
            for post in example_posts:
                if post['categoryId'] == cg['id']:
                    cg['horizontalcardproductimageurl'] = post['horizontalcardproductimageurl']
                    sub_cats.append(cg)
                    break
        ctgs.append(sub_cats)

    parent_group = define_silo(cats, 'category')

    google_analytics_id = os.environ.get("GOOGLE_ANALYTICS_ID")

    context = {
        "categories": ctgs,
        "parent_group": parent_group,
        "current_year": current_year,
        "example_posts": example_posts[:10],
        "domain": domain,
        "site_name": site_name,
        "logo": logo,
        "page_title": page_title,
        "meta_description": meta_description,
        "google_analytics_id": google_analytics_id,
        "google_analytics_src": "https://www.googletagmanager.com/gtag/js?id=" + google_analytics_id
    }

    return render(request, 'blog/category_group.html', context)


def category(request, *args, **kwargs):
    category_slug = kwargs['slug']

    cursor = connection.cursor()
    cursor.execute(
        f'''
        SELECT JSON_BUILD_OBJECT('posts', (
            SELECT JSON_AGG(row_to_json(review_post)) FROM review_post WHERE review_post."categoryId" = (
                SELECT id FROM category AS c WHERE c.slug = '{category_slug}'
            )
        ), 'parent_group', (
            SELECT JSON_AGG(TO_JSON(parent_group)) AS p FROM parent_group
        ), 'category_group', (
            SELECT JSON_AGG(TO_JSON(category_group)) AS p FROM category_group
        )) AS category;
        '''
    )
    rows = cursor.fetchall()
    desc = cursor.description
    columns = [col[0] for col in desc]
    results = [
        dict(zip(columns, row))
        for row in rows
    ]

    if not results[0]['category']['posts']:
        return HttpResponseNotFound("Not found.")

    path = request.path.split('/category/')[1].replace("-", " ")
    page_title = "Reviews of " + path.title() + " | " + site_name
    meta_description = f"Learn more about the {path} with the {site_name}."

    posts = results[0]['category']['posts']

    parent_group = define_silo(query_result=results, key_name='category')

    separated_posts = loop_items(items=posts)

    google_analytics_id = os.environ.get("GOOGLE_ANALYTICS_ID")

    context = {
        "posts": separated_posts,
        "domain": domain,
        "site_name": site_name,
        "logo": logo,
        "page_title": page_title,
        "parent_group": parent_group,
        "path": request.path,
        "meta_description": meta_description,
        "current_year": current_year,
        "google_analytics_id": google_analytics_id,
        "google_analytics_src": "https://www.googletagmanager.com/gtag/js?id=" + google_analytics_id
    }

    return render(request, 'blog/category.html', context)


def review_post(request, *args, **kwargs):
    rp_slug = kwargs['slug']

    # Get Post Data
    cursor = connection.cursor()
    cursor.execute(
        f'''
        SELECT JSON_BUILD_OBJECT('post', (
            SELECT JSON_AGG(row_to_json(review_post)) AS post FROM review_post WHERE slug = '{rp_slug}'
        ), 'related_rev', (
            SELECT JSON_AGG(TO_JSON(review_post)) AS related_rev FROM review_post WHERE "categoryId" = (
                SELECT "categoryId" FROM review_post AS rp WHERE rp.slug = '{rp_slug}'
            ) GROUP BY "categoryId"
        ), 'product_info', (
            SELECT JSON_AGG(TO_JSON(product)) AS p FROM product WHERE "affiliateUrl" = (
                SELECT productaffiliateurl FROM review_post AS rp WHERE rp.slug = '{rp_slug}'
            )
        ), 'category', (
            SELECT JSON_AGG(TO_JSON(category)) AS p FROM category WHERE id = (
                SELECT "categoryId" FROM review_post AS rp WHERE rp.slug = '{rp_slug}'
            )
        ), 'parent_group', (
            SELECT JSON_AGG(TO_JSON(parent_group)) AS p FROM parent_group
        ), 'category_group', (
            SELECT JSON_AGG(TO_JSON(category_group)) AS p FROM category_group
        )) AS review_post;
        '''
    )
    rows = cursor.fetchall()
    desc = cursor.description
    columns = [col[0] for col in desc]
    rev_post = [
        dict(zip(columns, row))
        for row in rows
    ]

    if not rev_post[0]['review_post']['post']:
        return HttpResponseNotFound("Not found.")

    # Get Product Information
    product_info = rev_post[0]['review_post']['product_info'][0]

    # Output Stars
    i = 0
    product_rating_stars = ""
    if product_info['productRatings'] == "":
        product_rating_stars = 0
    else:
        while i < float(product_info['productRatings']):
            product_rating_stars += '<li><i class="flaticon-star"></i></li>'
            i += 1

    post = rev_post[0]['review_post']['post'][0]
    related_rev = rev_post[0]['review_post']['related_rev']

    parent_group = define_silo(query_result=rev_post, key_name='review_post')
    page_title = post['title']
    meta_description = post['description']

    google_analytics_id = os.environ.get("GOOGLE_ANALYTICS_ID")

    context = {
        "post": post,
        "product_info": product_info,
        "related_rev": related_rev,
        "domain": domain,
        "site_name": site_name,
        "logo": logo,
        "product_rating_stars": product_rating_stars,
        "category": rev_post[0]['review_post']['category'][0],
        "parent_group": parent_group,
        "path": request.path,
        "meta_description": meta_description,
        "page_title": page_title,
        "current_year": current_year,
        "google_analytics_id": google_analytics_id,
        "google_analytics_src": "https://www.googletagmanager.com/gtag/js?id=" + google_analytics_id
    }

    return render(request, 'blog/review_post.html', context)


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
        "current_year": current_year
    }
    return render(request, 'blog/sitemap.xml.gz', context, content_type="application/xhtml+xml")


def sitemap_index(request, *args, **kwargs):

    cursor = connection.cursor()
    cursor.execute(
        '''SELECT CONCAT('0.8'), slug FROM review_post;''')
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
        "current_year": current_year
    }
    return render(request, 'blog/sitemap_index.xml', context, content_type="application/xhtml+xml")


def affiliate_disclaimer(request, *args, **kwargs):
    # Get Post Data
    cursor = connection.cursor()
    cursor.execute(
        f'''
        SELECT JSON_BUILD_OBJECT('parent_group', (
            SELECT JSON_AGG(TO_JSON(parent_group)) AS p FROM parent_group
        ), 'category_group', (
            SELECT JSON_AGG(TO_JSON(category_group)) AS p FROM category_group
        )) AS parent_group;
        '''
    )
    rows = cursor.fetchall()
    desc = cursor.description
    columns = [col[0] for col in desc]
    data = [
        dict(zip(columns, row))
        for row in rows
    ]

    parent_group = define_silo(query_result=data, key_name='parent_group')

    google_analytics_id = os.environ.get("GOOGLE_ANALYTICS_ID")

    context = {
        "parent_group": parent_group,
        "site_name": site_name,
        "current_year": current_year,
        "google_analytics_id": google_analytics_id,
        "google_analytics_src": "https://www.googletagmanager.com/gtag/js?id=" + google_analytics_id
    }

    return render(request, 'blog/disclaimer.html', context)


def privacy_policy(request, *args, **kwargs):
    # Get Post Data
    cursor = connection.cursor()
    cursor.execute(
        f'''
        SELECT JSON_BUILD_OBJECT('parent_group', (
            SELECT JSON_AGG(TO_JSON(parent_group)) AS p FROM parent_group
        ), 'category_group', (
            SELECT JSON_AGG(TO_JSON(category_group)) AS p FROM category_group
        )) AS parent_group;
        '''
    )
    rows = cursor.fetchall()
    desc = cursor.description
    columns = [col[0] for col in desc]
    data = [
        dict(zip(columns, row))
        for row in rows
    ]

    parent_group = define_silo(query_result=data, key_name='parent_group')

    google_analytics_id = os.environ.get("GOOGLE_ANALYTICS_ID")

    context = {
        "parent_group": parent_group,
        "site_name": site_name,
        "current_year": current_year,
        "domain": domain,
        "google_analytics_id": google_analytics_id,
        "google_analytics_src": "https://www.googletagmanager.com/gtag/js?id=" + google_analytics_id
    }

    return render(request, 'blog/privacy_policy.html', context)
