from django.urls import path
from . import views

urlpatterns = [
    path('', views.home, name='home'),
    path('<slug:sub_category_slug>/<slug:slug>/', views.ReviewPostView.as_view(), name='review_post'),
    path('<slug:category_slug>/<slug:sub_category_slug>/', views.SubCategoryView.as_view(), name='sub_category'),
    path('<slug:category_slug>/', views.CategoryView.as_view(), name='category')
]