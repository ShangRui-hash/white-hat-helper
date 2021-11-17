# Create your views here.
#coding=utf-8

from django.shortcuts import render
from django import forms
from django.http import HttpResponse
from django.contrib import auth
from account.models import User
from lib.utils import response

class UserFormRegister(forms.Form):
    username = forms.CharField(label='用户名:',max_length=100)
    password1 = forms.CharField(label='密码:',widget=forms.PasswordInput())
    password2 = forms.CharField(label='确认密码:',widget=forms.PasswordInput())
    email = forms.EmailField(label='电子邮件:')

class UserFormLogin(forms.Form):
    username = forms.CharField(label='用户名:',max_length=100)
    password = forms.CharField(label='密码:',widget=forms.PasswordInput())


def index(request):
    pass

def login(request):
    if request.method == "POST":
        uf = UserFormLogin(request.POST)
        if uf.is_valid():
        
            username = uf.cleaned_data['username']
            password = uf.cleaned_data['password']  
            userResult = User.objects.filter(username=username,password=password)
            
            if (len(userResult)>0):
                return response.resp_success("登录成功")
            
            else:
              return response.resp_fail("该用户不存在")
                #返回登陆页面
    
       

def register(request):
    
    if request.method == "POST":
        uf = UserFormRegister(request.POST)
        if uf.is_valid():
            #获取表单信息
            username = uf.cleaned_data['username']
            filterResult = User.objects.filter(username = username)
            if len(filterResult)>0:
                
                return response.resp_fail("用户名已存在")

            else:
                password1 = uf.cleaned_data['password1']
                password2 = uf.cleaned_data['password2']
                
                if (password2 != password1):
                   
                    return response.resp_fail("两次输入的密码不一致!")
                    #返回注册页面
                else:
                    password = password2
                    email = uf.cleaned_data['email']
                    #将表单写入数据库
                    user=User()
                    user.name = username
                    user.password = password1
                    user.email = email              
                    user.save()
                    #注册成功
                    return response.resp_success("注册成功")

    
      
