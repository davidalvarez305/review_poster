from django.urls import path
from . import views

urlpatterns = [
    path('', views.home, name='home'),
    path('category/<slug:slug>', views.category, name='category'),
    path('single/<slug:slug>', views.review_post, name='review_post'),
    path('category-group/<slug:name>', views.category_group, name='category_group'),
    path('<int:int>/sitemap.xml.gz', views.sitemap, name='sitemap'),
    path('sitemap_index.xml', views.sitemap_index, name='sitemap_index'),
    path('disclaimer', views.affiliate_disclaimer, name='disclaimer'),
    path('privacy-policy', views.privacy_policy, name='privacy_policy')
]
