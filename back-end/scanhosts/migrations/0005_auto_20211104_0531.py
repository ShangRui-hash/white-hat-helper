# Generated by Django 3.2.8 on 2021-11-04 05:31

from django.db import migrations


class Migration(migrations.Migration):

    dependencies = [
        ('scanhosts', '0004_task'),
    ]

    operations = [
        migrations.DeleteModel(
            name='BrowserInfo',
        ),
        migrations.DeleteModel(
            name='UserIPInfo',
        ),
    ]
