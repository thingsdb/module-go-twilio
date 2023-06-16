# SMTP ThingsDB Module (Go)

SMTP module written using the [Go language](https://golang.org).


## Installation

Install the module by running the following command in the `@thingsdb` scope:

```javascript
new_module("smtp", "github.com/thingsdb/module-go-smtp");
```

Optionally, you can choose a specific version by adding a `@` followed with the release tag. For example: `@v0.1.0`.

## Configuration

The smtp module requires configuration with the following properties:

Property | Type            | Description
-------- | --------------- | -----------
host     | str (required)  | SMTP host, eg 'myhost.local:587'
auth     | [str, str]      | Optional authentication. [Username, Password].


Example configuration:

```javascript
set_module_conf("twilio", {
    host: "myhost.local:587",
    auth: ["myuser", "mypassword"],
});
```

## Exposed functions

Name                            | Description
------------------------------- | -----------
[send_mail](#send-mail)         | Send an email.

### Send mail

Syntax: `send_mail(to, mail)`

#### Arguments

- `mail`: _(thing)_ Mail to send.

#### Example:

```javascript
// Only subject is required
mail = {
    bcc: ['charlie@foo.bar'],
    cc: ['info@foo.bar'],
    from: 'bob@foo.bar',
    from_name: 'Bob',
    html: '<html>Html Body</html>',
    plain: 'plain text body',
    reply_to: 'bob@foo.bar',
    subject: 'my subject',
};

// At least one to address is required
to = ['alice@foo.bar'];

// Send the email
smtp.send_mail(to, mail).else(|err| {
    err;  // some error has occurred
})
```

# module-go-twilio
