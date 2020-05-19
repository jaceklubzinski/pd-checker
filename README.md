# pd-checker-service


# pd-checker-event
## Description
Main idea behind this program is to test _PagerDuty_ integration on your infrastracture.

pd-checker-event trigger (from inside of your infrastracture) and instantly resolve single Pagerduty incident always with the same options:
```
Summary:  "PD CHECKER - OK",
Severity: "info",
Source:   "localhost",
Details:  triggerEvery,
```
New event can be create manually or in server mode every _triggerEvery_ time.

Next pd-checker-service will scan all available services every _triggerEvery_ time and if found new event with name _PD CHECKER - OK_ register them in local database 


## Usage
Integration key _Events API v2_
test
