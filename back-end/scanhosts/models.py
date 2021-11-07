from django.db import models
import time

# Create your models here.
class Company(models.Model):
    """公司表"""
    name = models.CharField(max_length=30)
    created_at = models.DateTimeField(
        default=time.strftime("%Y-%m-%d %H:%M:%S", time.localtime()))
    updated_at = models.DateTimeField(auto_now=True)


class Host(models.Model):
    """主机表"""
    ip = models.CharField(max_length=15)
    os = models.CharField(max_length=30)
    domain_list = models.CharField(max_length=300, blank=True)
    ports_list = models.CharField(max_length=300, blank=True)
    created_at = models.DateTimeField(
        default=time.strftime("%Y-%m-%d %H:%M:%S", time.localtime()))
    company = models.ForeignKey("Company", on_delete=models.CASCADE)
    updated_at = models.DateTimeField(auto_now=True)


class Task(models.Model):
    """任务表"""
    targets = models.CharField(max_length=300)
    company = models.ForeignKey("Company", on_delete=models.CASCADE)
    status = models.CharField(max_length=30, default="pending")
    created_at = models.DateTimeField(
        default=time.strftime("%Y-%m-%d %H:%M:%S", time.localtime()))
    updated_at = models.DateTimeField(auto_now=True)
