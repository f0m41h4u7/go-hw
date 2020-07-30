# file: features/api.feature
# language: en-pirate

Ahoy matey!: Calendar API
  As API client of Calendar service
  In order to work with events
  I want to CRUD events via API

  Heave to: Create event
    Blimey! I send "POST" request to "http://calendar:1337/create" with body:
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
    Blimey! I send "GET" request to "http://calendar:1337/getday?day=2049-07-07T00:00:00"
    Let go and haul The response code should be 200
    Aye I receive event
