from datetime import date
import os
from django.shortcuts import render

# Environment Variables
DOMAIN = str(os.environ.get('DOMAIN'))
CURRENT_YEAR = date.today().year

def home(request):
    context = {
        "domain": DOMAIN,
        "current_year": CURRENT_YEAR
    }
    return render(request, 'blog/home.html', context)