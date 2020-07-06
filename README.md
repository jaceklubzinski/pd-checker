# Architecure overview
![Diagram](doc/images/pd-checker.png)

#pd-checker-service

- Get all services from PagerDuty and save them to sqlite database (services table)
- Check pd-checker event for services that don't have yet any event and save them to sqlite database (incidents table)
- Get all incidents from sqlite database
- Mark incident to check if created time is bigger than defined `triggerEvery` for this incident

### Database structure
`pd-checker-service` lists incidents for all available services and save to the database the last incident per service created by the `pd-checker-event`.
If incident for given service already exist it will only update `ID`, `Title`, `CreateAt` and `Timer` values. 
```
//Incident structure for incidents stored in database
type Incident struct {
	ID          string // PagerDuty incident ID
	Title       string // PagerDuty incident title
	ServiceID   string // PagerDuty service ID related to created incident
	ServiceName string // PagerDuty service name related to created incident
	CreateAt    string // PagerDuty incident creation time
	Timer       string // PagerDuty additional information defined by pd-checker event details informtion
	Alert       string // If "Y" create new alert for service
	ToCheck     string // If "Y" check for new incidents 
	Trigger     string // If "Y" alert already triggered
}
```

All available PagerDuty services are stored in `Service` database.
```
//Service structure for service stored in database
type Service struct {
	ID   string // PagerDuty service ID related to created incident
	Name string // PagerDuty service name related to created incident
}
```

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
