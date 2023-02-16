from django.urls import path
from . import views

urlpatterns = [
    path('', views.home, name='home'),
    path('<slug:sub_category>/<slug:slug>/', views.review_post, name='review_post')
]