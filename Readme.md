# Endpoints:
> ## Auth:
> > #### Register: /auth/sign-up
> > #### Authorization: /auth/sign-in
> ## Categories:
> > #### Create Category: /seller/category
> > #### Update Category: /seller/category/:category_id
> > #### Delete Category: /seller/category/:category_id
> > #### Get Category: /seller/category/:category_id
> > #### Get All Categories: /seller/category/
> ## Subcategories:
> > #### Create Subcategory: /seller/category/:category_id
> > #### Update Subcategory: /seller/category/:cat_id/subcategory/:subcat_id
> > #### Delete Subcategory: /seller/category/:cat_id/subcategory/:subcat_id
> > #### Get Subcategory: /seller/category/:cat_id/subcategory/:subcat_id
> ## Products:
> > #### Create One Product: /seller/category/:cat_id/subcategory/:subcat_id/one
> > #### Create Many Products: /seller/category/:cat_id/subcategory/:subcat_id/many
> > ### Update Product: /seller/category/:cat_id/subcategory/:subcat_id/products/:prod_id
> > ### Delete Product: /seller/category/:cat_id/subcategory/:subcat_id/products/:prod_id
> ## History:
> > #### Get History: /seller/history
> > #### Get Transaction: /seller/history/:tran_id
> ## Client:
> > #### Get Product: /api/client/:username?uniquecode=....
> > #### Pre check: /api/precheck