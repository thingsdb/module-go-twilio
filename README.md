# Twilio ThingsDB Module (Go)

Twilio module written using the [Go language](https://golang.org).


## Installation

Install the module by running the following command in the `@thingsdb` scope:

```javascript
new_module("twilio", "github.com/thingsdb/module-go-twilio");
```

Optionally, you can choose a specific version by adding a `@` followed with the release tag. For example: `@v0.1.1`.

## Configuration

The Twilio module requires configuration with the following properties:

Property             | Type            | Description
-------------------- | --------------- | -----------
`TWILIO_ACCOUNT_SID` | str (required)  | Find your Account SID at twilio.com/console.
`TWILIO_AUTH_TOKEN`  | str (required)  | Find your Auth Token twilio.com/console.


Example configuration:

```javascript
set_module_conf("twilio", {
    TWILIO_ACCOUNT_SID: "REPLACE WITH ACCOUNT_SID",
    TWILIO_AUTH_TOKEN: "REPLACE WITH AUTH_TOKEN",
});
```

## Exposed functions

Name                                    | Description
--------------------------------------- | -----------
[call](#voice-call)                     | Make a Voice call.
[message (SMS)](#sms-message)           | Sent a SMS message
[message (WhatsApp)](#whatsapp-message) | Sent a WhatsApp message

### Voice call

Syntax: `call(params)`

#### Arguments

- `params`: _(thing)_ Params to use for the voice call.

#### Example:

```javascript
// Only subject is required
params = {
    body: 'Hello world!',
    to: '+310612345678',
    from: '+310687654321',
};

// Make the call
twilio.call(params).then(|resp| {
    // you might want to do something with the response, for example resp.Sid.
}).else(|err| {
    err;  // some error has occurred
});
```

### SMS Message

Syntax: `message(params)`

#### Arguments

- `params`: _(thing)_ Params to use for the voice call.

#### Example:

```javascript
// Only subject is required
params = {
    body: 'Hello world!',
    to: '+310612345678',
    from: '+310687654321',
};

// Make the call
twilio.message(params).then(|resp| {
    // you might want to do something with the response, for example resp.Sid.
}).else(|err| {
    err;  // some error has occurred
});
```

### WhatsApp Message

Syntax: `message(params)`

#### Arguments

- `params`: _(thing)_ Params to use for the voice call.

#### Example:

```javascript
// Only subject is required
params = {
    body: 'Hello world!',
    to: 'whatsapp:+310612345678',
    from: 'whatsapp:+310687654321',
};

// Make the call
twilio.message(params).then(|resp| {
    // you might want to do something with the response, for example resp.Sid.
}).else(|err| {
    err;  // some error has occurred
});
```
