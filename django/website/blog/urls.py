from django.urls import path
from . import views

urlpatterns = [
    path('', views.HomeView.as_view(), name='home'),
    path('<slug:category>/<slug:sub_category>/<slug:slug>', views.ReviewPostView.as_view(), name='review_post'),
    path('<slug:category>/<slug:sub_category>/', views.SubCategoryView.as_view(), name='sub_category'),
    path('<slug:category>/', views.CategoryView.as_view(), name='category'),
    path('<int:int>/sitemap.xml.gz', views.sitemap, name='sitemap'),
    path('sitemap_index.xml', views.sitemap_index, name='sitemap_index'),
    path('disclaimer', views.AffiliateDisclaimer.as_view(), name='disclaimer'),
    path('privacy-policy', views.PrivacyPolicy.as_view(), name='privacy_policy'),
    path('create-review-post', views.CreatePost.as_view(), name='create_review_post'),
    path('login', views.Login.as_view(), name='login'),
    path('logout', views.Logout.as_view(), name='logout')
]

#path('create-post', views.CreatePost.as_view(), name='create_post')