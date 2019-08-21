Viridian chaincode specifications
=================================

This file is about the code. See `model` and `output` for the data model that the code operates on.

This file defines the behavior of the chaincode API, i.e. the functions that can be used to interact with the blockchain. It is meant for TDD (test-driven development) and BDD (behavior-driven development) purposes.

What functionality must the functions provide? What is supposed to happen when a function is given a certain input? What kinds of wrong usage can there be and how is the code supposed to handle it?

Based on the specifications in this file, tests can be written that will ensure that the code is behaving properly.

Product specification
---------------------

\* = required<br>
<sup>&dagger;</sup> = must be unique

* *addProduct:* It should be possible to add a new product.
    * **Inputs:**
        * Product key (uuid)\*<sup>&dagger;</sup>,
        * GTIN
        * Producer key
        * Contained product keys
        * Label keys
        * Locales (at least one structure with the following fields):
            * lang\*<sup>&dagger;</sup>
            * name\*
            * price
            * currecny
            * description
            * quantities
            * ingredients
            * packagings
            * categories
            * imageURL
            * URL
    * **Results/Side Effects:**
        * A new product is registered under the key, it should have status "Preliminary" until review closed (either passed or not passed)
        * A review should be created and random users assigned to it
        * If review passed:
            * The product should now have status "Active"
        * If review not passed:
            * The product should now have status "Rejected"
        * Ideally, when the review is created, it is searched for similar existing products and they are stored as "Possible duplicate of ..." with the review.
        * Even more ideally, the similar products are already displayed to the user during add process with the ability to modify one of those existing products
    * **Edge Cases:**
        * No product key provided
        * Product key already used
        * Producer key not found in blockchain
        * Contained product keys not found in blockchain
        * Label keys not found in blockchain
        * Not even one locale
        * More than one locale with same lang
        * Also add regex checks for GTIN, lang, price, currency, URLs etc.?
        * Submitting user not registered
* *editProduct:* It should be possible to edit (i.e. modify) a product, but only if its status is "Active" and there is no edit/deletion pending (=in the review queue, which is signified by `supersededBy` not being empty).
    * **Inputs:**
        * Old product key\*
        * Change reason\* (really make this required?)
        * Rest is same as "addProduct":
            * New product key (uuid)\*<sup>&dagger;</sup>
            * GTIN
            * Producer key
            * Contained product keys
            * Label keys
            * Locales (at least one structure with the following fields):
                * lang\*<sup>&dagger;</sup>
                * name\*
                * price
                * currecny
                * description
                * quantities
                * ingredients
                * packagings
                * categories
                * imageURL
                * URL
    * **Results/Side Effects:**
        * A new product is registered under the new key, it should have the old product key under `supersedes`
        * The old product should have the new product key under `supersededBy`
        * A review should be created and random users assigned to it
        * The new product should have status "Preliminary" until review closed
        * The old product should continue to have status "Active" until review closed (either passed or not passed)
        * If review passed:
            * The old product should now have status "Outdated"
            * The new product should now have status "Active"
        * If review not passed:
            * The old product should continue to have status "Active"
            * The new product should now have status "Rejected"
            * The old product's `supersededBy` should again be empty
        * Ideally, when the review is created, it is searched for similar existing products and they are stored as "Possible duplicate of ..." with the review.
        * Even more ideally, the similar products are already displayed to the user during add process with the ability to modify one of those existing products
    * **Edge Cases:**
        * Old product key not found
        * Old product does not have status "Active"
        * Old product's `supersededBy` is not empty
        * No change reason provided (really make this required?)
        * Rest is same as "addProduct":
            * No new product key provided
            * New product key already used
            * Producer key not found in blockchain
            * Contained product keys not found in blockchain
            * Label keys not found in blockchain
            * Not even one locale
            * More than one locale with same lang
            * Also add regex checks for GTIN, lang, price, currency, URLs etc.?
            * Submitting user not registered
* *deleteProduct:* It should be possible to delete (i.e. remove) a product, but only if its status is "Active" and there is no edit/deletion pending (=in the review queue, which is signified by `supersededBy` not being empty).
    * **Inputs:**
        * Product key\*
        * Change, i.e. delete reason\* (really make this required?)
    * **Results/Side Effects:**
        * A review should be created and random users assigned to it
        * The product should continue to have status "Active" until review closed (either passed or not passed)
        * The product should have `supersededBy` set to the string "DELETION"
        * If review passed:
            * The product should now have status "Deleted"
        * If review not passed:
            * The product should continue to have status "Active"
            * The product's `supersededBy` should again be empty
    * **Edge Cases:**
        * Product key not found
        * Product does not have status "Active"
        * Product's `supersededBy` is not empty
        * No change reason provided (really make this required?)
        * Submitting user not registered