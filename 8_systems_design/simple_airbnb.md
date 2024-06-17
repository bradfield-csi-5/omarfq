# Simple AirBnbB

## Requirements Document
[https://docs.google.com/document/d/1JSi1MHRGs69XnmgkzLpuZuWaO1UZo-NBKP50lv6-qww/edit#heading=h.o1usn3j8q9ik]()

## Notes from video:
- People should be able to book a room transactionally, meaning two people can't book the same room at the same time.
- Let's assume that, for location purposes, they will zoom in/out a map. If we'd like the users to search by city, let's also assume that we can reach out to some other place, find out where New York City is (for example), and then zoom them in on that map. The map becomes the interface for the user. Assume we have some way to match cities with rectangles on the map.
  - For this, assume that we can call Google Maps' API and obtain the coordinates for a given city.
- The reservation/booking cannot be modified once it's made.
- Our design should be extensible, meaning that if we wanted to integrate a system for reviews we should be able to do it seamlessly. We shouldn't worry about reviews.
- Assume there are transactional messages (e.g., an email arrives when a booking is complete).
- Phone numbers can be used to communicate the client with the host.
- We don't need to implement an auth system, the user can already sign in with Google, Facebook, etc.
- The host can't reject a booking. If a room is available and someone books then that's final.

## Additional Assumptions
- 1 million places, each of them has pictures and a description.
  - Assume the host can upload 5 images of the room/property.
    - Each image has a size of say 2 MB.
  - Assume the description has a max length of 2000 characters.
- On data access: assume a Pareto distribution on users browsing vs users booking. Also, number of potential guests is significantly bigger than the number of hosts.
- Assume that if a user filters by a city in the map view, our system will make a request to Google Maps API to get the longitude and latitude for, let's say, New York. Once our system receives those coordinates, it will zoom in on the map the place will show up within the geo-rectangle for this city. Does this mean we don't need to store anything related to location in our models...?

## User Journeys
- A user can sign up as either a guest or host.
- A guest books a place.
- A host posts a place for booking.

## Data Models
- `user`
  - user_id - BIGINT -> 8 bytes (auto-increment) PK
  - name - VARCHAR(255) -> 255 * 4 = 1020 bytes
  - last_name - VARCHAR(255) -> 255 * 4 = 1020 bytes
  - email - VARCHAR(255) -> 255 * 4 = 1020 bytes
  - phone_number - VARCHAR(20) -> 20 * 4 = 80 bytes
  - role - ENUM('guest', host') -> 4 bytes
- `listing`
  - listing_id - BIGINT -> 8 bytes (auto-increment) PK
  - user_id - BIGINT -> 8 bytes FK reference to `user`
  - description - VARCHAR(2000) -> 2000 * 4 = 8000 bytes
  - price - INT -> 4 bytes
  - available_from - TIMESTAMP -> 8 bytes
  - available_until - TIMESTAMP -> 8 bytes
  - CONSTRAINT available_until > available_from
- `listing_image`
  - image_id - BIGINT -> 8 bytes (auto-increment) PK
  - listing_id - BIGINT -> 8 bytes FK reference to `listing`
  - image_url - TEXT -> assume 2000 bytes
- `booking`
  - booking_id - BIGINT -> 8 bytes (auto-increment) PK
  - user_id - BIGINT -> 8 bytes
  - listing_id - BIGINT -> 8 bytes
  - happened_at - TIMESTAMP -> 8 bytes
  - expires_at - TIMESTAMP -> 8 bytes
  - status ENUM('free', 'booked', 'reserved', 'expired') -> 4 bytes
 
## Endpoints
- GET /listings -> Returns all available listings pag

## Calculations


## What high-level components do we need?
- A frontend where the user can.
- A server to process requests sent by the frontend.
- A database where we can store user info, listings, perform transactions, etc.
- We may need to scale these components appropriately.
