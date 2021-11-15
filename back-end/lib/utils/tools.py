from django.core.mail import send_mail
import time

class sendmail():
    def __init__(self,receiver,subject,content):
        self.receiver = receiver
        self.sender = "shangrui1024@163.com"
        self.subject = subject + ' ' + time.strftime('%Y-%m-%d %H:%M:%S',time.localtime(time.time()))
        self.content = content
        
    def send(self):
        try:
            send_mail(self.subject,self.content,self.sender,self.receiver,fail_silently=False)
            return True 
        except Exception as e:
            print("sendmail failed,err:",e)
            return False