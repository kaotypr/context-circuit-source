---
kind: role
title: Shopper
owners:
  - Storefront product team
sources:
  - Acme Storefront PRD v3
review_date: 2026-01-05
relevant_domains:
  - ../domains/checkout/README.md
related_workflows:
  - ../domains/checkout/workflows/place-order.md
known_gaps:
  - Guest vs registered checkout differences are not documented.
---

# Shopper

## Role definition

A customer who browses products and buys them through the storefront.

## Primary outcomes

Find a product, pay for it, and receive an order confirmation.

## Product surfaces

The storefront web app: catalog, cart, and checkout screens.

## End-to-end role story

A shopper adds products to a cart and completes [placing an order](../domains/checkout/workflows/place-order.md)
to receive a confirmation. Exact payment rules live in that workflow page.

## Related workflows

- [Place an order](../domains/checkout/workflows/place-order.md)

## Role-specific behavior

Shoppers can only see their own orders.

## Limitations

A shopper cannot edit an order once it is paid.
