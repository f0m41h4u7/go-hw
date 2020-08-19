# language: en-pirate

Ahoy matey!: Calendar API
  As API client of Calendar service
  In order to work with events
  I want to CRUD events via API

  Yo-ho-ho: Calendar address is "http://calendar:1337"

  Heave to: Create event
    Blimey! I send "POST" request to "/create" with body:
		"""
		{
			"title": "test event",
      "start": "2049-07-07T10:00:00",
      "end": "2049-07-07T12:00:00",
      "description": "some test event",
      "owner_id": "9bed7c53-c3bd-4f7e-92d1-5d98c04fb83a",
      "notify_in": "2h"
		}
		"""
    Let go and haul The response code should be 200
		Aye I receive UUID

  Heave to: Get event
    Blimey! I send "GET" request to "/getday?day=2049-07-07T00:00:00"
    Let go and haul The response code should be 200
    Aye I receive event

  Heave to: Update event
    Blimey! I send "PUT" request to "/update" updating title to "updated event"
    Let go and haul The response code should be 200
    Let go and haul I send "GET" request to "/getday?day=2049-07-07T00:00:00"
    Aye I receive updated event

  Heave to: Get events for week
    Gangway! There are events in DB
    Blimey! I send "GET" request to "/getweek?week=2049-07-07T00:00:00"
    Let go and haul The response code should be 200
    Aye I receive list of events for week

  Heave to: Get events for month
    Blimey! I send "GET" request to "/getmonth?month=2049-07-07T00:00:00"
    Let go and haul The response code should be 200
    Aye I receive list of events for month

  Heave to: Delete event
    Blimey! I send "DELETE" request to "/delete"
    Let go and haul The response code should be 200
    Aye I receive no events at "/getday?day=2049-07-07T00:00:00"
