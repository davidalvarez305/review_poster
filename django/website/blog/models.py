from django.db import models

class Categorization(models.Model):
    name = models.CharField(max_length=250)
    slug = models.SlugField(db_index=True, unique=True)

    def __str__(self):
        return self.name

    class Meta:
        abstract = True

class Group(Categorization):

    class Meta:
        db_table = "group"

    def prefetch_category_set():
        groups = Group.objects.all().prefetch_related('category_set').all()

        return groups

class Category(Categorization):
    group = models.ForeignKey(Group, on_delete=models.CASCADE)

    def __str__(self):
        return self.name

    class Meta:
        db_table = "category"

class SubCategory(Categorization):
    category = models.ForeignKey(Category, on_delete=models.CASCADE)

    def __str__(self):
        return self.name

    class Meta:
        db_table = "sub_category"

class Product(models.Model):
    affiliate_url = models.CharField(max_length=250, db_index=True, unique=True)
    product_price = models.CharField(max_length=250, default='')
    product_reviews = models.CharField(max_length=250, default='')
    product_ratings = models.CharField(max_length=250, default='')
    product_image = models.CharField(max_length=250, default='')
    product_image_alt = models.CharField(max_length=250, default='')
    product_label = models.CharField(max_length=250, default='')
    product_name = models.CharField(max_length=250, default='')
    product_description = models.TextField(default='')

    def __str__(self):
        return self.affiliate_url

    class Meta:
        db_table = "product"

class ReviewPost(models.Model):
    title = models.CharField(max_length=250)
    slug = models.SlugField(db_index=True, unique=True)
    content = models.TextField()
    sub_category = models.ForeignKey(SubCategory, on_delete=models.CASCADE)
    headline = models.TextField()
    intro = models.TextField()
    description = models.TextField()
    product_affiliate_url = models.ForeignKey(Product, db_column='product_affiliate_url', to_field='affiliate_url', on_delete=models.CASCADE)
    faq_answer_1 = models.TextField()
    faq_answer_2 = models.TextField()
    faq_answer_3 = models.TextField()
    faq_question_1 = models.TextField()
    faq_question_2 = models.TextField()
    faq_question_3 = models.TextField()

    def __str__(self):
        return self.title

    class Meta:
        db_table = "review_post"