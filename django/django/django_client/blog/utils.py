import math


def define_silo(query_result, key_name):
    parent_group = query_result[0][key_name]['parent_group']
    category_group = query_result[0][key_name]['category_group']

    for index, pg in enumerate(parent_group):
        related_categories = []
        sub_list = []
        count = 0
        for cg in category_group:
            if cg['parent_group_id'] == pg['id']:
                if count == 5:
                    related_categories.append(sub_list)
                    count = 0
                    sub_list = []
                sub_list.append(cg)
                count += 1

        pg['category_group'] = related_categories
        parent_group[index] = pg
    return parent_group


def loop_items(items):
    ctgs = []
    sub_list = []
    count = 0
    max_count = math.ceil(len(items) / 3)
    for cat in items:
        if count == max_count:
            ctgs.append(sub_list)
            count = 0
            sub_list = []
        sub_list.append(cat)
        count += 1

    return ctgs


def fix_images(posts):
    print(posts)
    fixed_related = []
    for related_p in posts:
        fixed_rp = related_p['horizontalcardproductimageurl'].replace(
            "cdn.southfloridaathleticclub.com", "southfloridaathleticclub.s3.amazonaws.com")
        related_p['horizontalcardproductimageurl'] = fixed_rp
        fixed_related.append(related_p)
    return fixed_related
