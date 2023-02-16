from django.db import models

class Category(models.Model):
    name = models.TextField()
    slug = models.SlugField(db_index=True)

    def __str__(self):
        return self.name

class SubCategory(models.Model):
    name = models.TextField()
    slug = models.SlugField(db_index=True)
    category_id = models.ForeignKey(Category, on_delete=models.CASCADE)

    def __str__(self):
        return self.name

class ReviewPost(models.Model):
    title = models.CharField()
    slug = models.SlugField(db_index=True)
    content = models.TextField()
    sub_category_id = models.ForeignKey(SubCategory, on_delete=models.CASCADE)
    headline = models.TextField()
    intro = models.TextField()
    description = models.TextField()
    product_label = models.TextField()
    product_name = models.TextField()
    product_description = models.TextField()
    product_affiliate_url = models.TextField()
    faq_answer_1 = models.TextField()
    faq_answer_2 = models.TextField()
    faq_answer_3 = models.TextField()
    faq_question_1 = models.TextField()
    faq_question_2 = models.TextField()
    faq_question_3 = models.TextField()
    product_image_url = models.TextField()
    product_image_alt = models.TextField()

    def __str__(self):
        return self.title