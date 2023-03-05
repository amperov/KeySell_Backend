CREATE TABLE sellers (
                         id serial primary key ,
                         username text not null unique ,
                         pass text not null ,
                         seller_id int not null unique ,
                         seller_key text not null unique
);
CREATE TABLE category(
                         id serial primary key ,
                         title_ru text not null ,
                         title_eng text,
                         item_id int,
                         description text,
                         message_client text,
                         created_at text,
                         user_id int references sellers (id) ON DELETE CASCADE
);
CREATE TABLE subcategory (
                             id serial primary key ,
                             title_ru text not null ,
                             title_eng text not null ,
                             subitem_id int,
                             created_at text,
                             subtype_value int,
                             partial_values text,
                             is_composite boolean,
                             category_id int references category (id) ON DELETE CASCADE
);
CREATE TABLE products (
                          id serial primary key ,
                          content_key text not null,
                          created_at text,
                          subcategory_id int references subcategory (id) ON DELETE CASCADE
);

CREATE TABLE transactions(
                             id serial primary key ,
                             category_name text ,
                             subcategory_name text  ,
                             client_email text,
                             content_key text,
                             amount int,
                             profit int,
                             amount_usd int,
                             count int,
                             unique_inv int,
                             user_id int references sellers (id),
                             unique_code text unique,
                             date_check text,
                             date_delivery text,
                             date_confirmed text,
                             state text
);