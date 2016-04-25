Introduction
============

This document describes the data models and JSON structure of the ghenga API,
thus describing the public data structures saved in the database. All data in
this document is fake and were generated at random.

In ghenga, the following entities are managed:

 * Person
 * Account (e.g. a company or a union)
 * Event
 * Activity
 * Task
 * User

Several fields of the models are managed by ghenga and cannot be changed, this
includes the following fields:

 * `id`
 * `created_at`
 * `changed_at`

Person
======

A Person is a contact somewhere. The JSON document describing a Person is as
follows:

```json
{
  "id": 100,
  "name": "Nicolai Person",
  "title": "CEO",
  "department": "Management",
  "email_addresses": "marlene@kleiningerkneifel.org",
  "phone_numbers": [
    {
      "type": "work",
      "number": "+49 221 1231234"
    },
    {
      "type": "mobile",
      "number": "+49 157 123123123"
    },
    {
      "type": "fax",
      "number": "+49 157 123123125"
    },
    {
      "type": "other",
      "number": "+49 221 1231235"
    }
  ],
  "address": {
    "street": "Teststraße 23",
    "postal_code": "50023",
    "state": null,
    "city": "Köln",
    "country": "Germany"
  },
  "comment": "This is a comment",
  "account_id": 123,
  "changed_at": "2016-04-24T10:30:07+00:00",
  "created_at": "2016-04-24T10:30:07+00:00"
}
```

Unset fields are either specified with the `null` value (see the field `state`
of the address), or not present in the JSON document. The field `phone_numbers`
is returned as an empty list when no phone numbers are present in the database.
The field `account_id` is the ID of an Account, and may also be `null`.

The following fields not automatically managed by ghenga are required for the
object to be valid:

 * `name`

Account
=======

An Account may be a company, a union or something else. The JSON document
describing an Account is as follows:

```json
{
  "id": 123,
  "name": "Beispiel GmbH",
  "website": "https://www.example.com",
  "phone_numbers": [
    {
      "type": "switchboard",
      "number": "+49 221 1231234"
    },
    {
      "type": "fax",
      "number": "+49 221 1231235"
    }
  ],
  "billing_address": {
    "street": "Teststraße 24",
    "postal_code": "03030",
    "state": null,
    "city": "Berlin",
    "country": "Germany"
  },
  "physical_address": {
    "street": "Teststraße 24b",
    "postal_code": "03030",
    "city": "Berlin",
    "country": "Germany"
  },
  "changed_at": "2016-04-24T10:30:07+00:00",
  "created_at": "2016-04-24T10:30:07+00:00"
}
```

The following fields not automatically managed by ghenga are required for the
object to be valid:

 * `name`
