# Generated by Django 3.2.8 on 2021-11-07 11:53

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('scanhosts', '0006_auto_20211104_1516'),
    ]

    operations = [
        migrations.AlterField(
            model_name='company',
            name='created_at',
            field=models.DateTimeField(default='2021-11-07 11:53:58'),
        ),
        migrations.AlterField(
            model_name='host',
            name='created_at',
            field=models.DateTimeField(default='2021-11-07 11:53:58'),
        ),
        migrations.AlterField(
            model_name='task',
            name='created_at',
            field=models.DateTimeField(default='2021-11-07 11:53:58'),
        ),
    ]
