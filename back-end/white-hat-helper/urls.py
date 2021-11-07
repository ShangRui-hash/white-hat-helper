"""mydevops URL Configuration

The `urlpatterns` list routes URLs to views. For more information please see:
    https://docs.djangoproject.com/en/3.2/topics/http/urls/
Examples:
Function views
    1. Add an import:  from my_app import views
    2. Add a URL to urlpatterns:  path('', views.home, name='home')
Class-based views
    1. Add an import:  from other_app.views import Home
    2. Add a URL to urlpatterns:  path('', Home.as_view(), name='home')
Including another URLconf
    1. Import the include() function: from django.urls import include, path
    2. Add a URL to urlpatterns:  path('blog/', include('blog.urls'))
"""
from django.contrib import admin
from django.urls import path
from django.views.decorators.csrf import csrf_exempt
import scanhosts.views as scanhosts
from showinfo.views.company import company_handler
from showinfo.views.hosts import hosts_handler

urlpatterns = [
    path('admin/', admin.site.urls),
    path('company',csrf_exempt(company_handler)),
    path('task',csrf_exempt(scanhosts.task_handler)),
    path('run',csrf_exempt(scanhosts.run_handler)),
    path('hosts',csrf_exempt(hosts_handler))
    # path('login',csrf_exempt(scanhosts.views.))
]
