from django.db import models

class Category(models.Model):
    name = models.CharField(max_length=250)
    slug = models.SlugField(db_index=True)

    def __str__(self):
        return self.name

    class Meta:
        db_table = "category"

class SubCategory(models.Model):
    name = models.CharField(max_length=250)
    slug = models.SlugField(db_index=True)
    category = models.ForeignKey(Category, on_delete=models.CASCADE)

    def __str__(self):
        return self.name

    class Meta:
        db_table = "sub_category"

class ReviewPost(models.Model):
    title = models.CharField(max_length=250)
    slug = models.SlugField(db_index=True)
    content = models.TextField()
    sub_category = models.ForeignKey(SubCategory, on_delete=models.CASCADE)
    headline = models.TextField()
    intro = models.TextField()
    description = models.TextField()
    product_label = models.CharField(max_length=250)
    product_name = models.CharField(max_length=250)
    product_description = models.TextField()
    product_affiliate_url = models.CharField(max_length=250)
    faq_answer_1 = models.TextField()
    faq_answer_2 = models.TextField()
    faq_answer_3 = models.TextField()
    faq_question_1 = models.TextField()
    faq_question_2 = models.TextField()
    faq_question_3 = models.TextField()
    product_image_url = models.CharField(max_length=250)
    product_image_alt = models.CharField(max_length=250)

    def __str__(self):
        return self.title

    class Meta:
        db_table = "review_post"