from django.urls import path
from . import views
from django.views.generic.base import RedirectView

urlpatterns = [
    path('', views.HomeView.as_view(), name='home'),
    path('scholarship/', RedirectView.as_view(url='/', permanent=True), name='scholarship'),
    path('strength-training-for-beginners/', RedirectView.as_view(url='/', permanent=True), name='strength_training_for_beginners'),
    path('best-powerlifting-belts/', RedirectView.as_view(url='/', permanent=True), name='best_powerlifting_belts'),
    path('best-power-racks/', RedirectView.as_view(url='/', permanent=True), name='best_power_racks'),
    path('beginner-powerlifting-program/', RedirectView.as_view(url='/', permanent=True), name='beginner_powerlifting_program'),
    path('how-often-should-you-train-for-strength/', RedirectView.as_view(url='/', permanent=True), name='how_often_should_you_train_for_strength'),
    path('5x5-workout/', RedirectView.as_view(url='/', permanent=True), name='5x5_workout'),
    path('powerlifting-home-gym/', RedirectView.as_view(url='/', permanent=True), name='powerlifting_home_gym'),
    path('single/<str:slug>/', RedirectView.as_view(url='/', permanent=True), name='redirect_single'),
    path('multiple/<str:slug>/', RedirectView.as_view(url='/', permanent=True), name='redirect_multiple'),
    path('training/<str:slug>/', RedirectView.as_view(url='/', permanent=True), name='redirect_training'),
    path('<slug:category>/<slug:sub_category>/<slug:slug>', views.ReviewPostView.as_view(), name='review_post'),
    path('<slug:category>/<slug:sub_category>/', views.SubCategoryView.as_view(), name='sub_category'),
    path('<slug:category>/', views.CategoryView.as_view(), name='category'),
    path('<int:int>/sitemap.xml.gz', views.sitemap, name='sitemap'),
    path('sitemap_index.xml', views.sitemap_index, name='sitemap_index'),
    path('disclaimer', views.AffiliateDisclaimer.as_view(), name='disclaimer'),
    path('privacy-policy', views.PrivacyPolicy.as_view(), name='privacy_policy')
]