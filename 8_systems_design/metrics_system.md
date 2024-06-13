# Metric System

## Initial Requirements
- Relatively larger company with tens of thousands of hosts
- We want to see how services and machines are performing
- What metrics are we concerned about:
  - Monitoring the health of our systems:
    - CPU (at this second, the machine is using 80% of CPU capacity), memory, network I/O, disk I/O
    - GC latency, out-of-memory errors, number of threads
    - How long are requests taking, how often are requests failing
- At what interval are we pulling system statistics
  - Every second is OK (every 10sec probably still OK)
- On latency, if something fails, users should be alerted within 10 secs of the failure ocurring so that they get to solve the issue quickly
  - If there's like a 2min lag users will be unhappy since its a pretty urgent situation
- Out of scope: logging errors / exceptions
- Consistency requirements:
  - It doesn't have to be perfect, but it shouldn't be "biased" (like you shouldn't get a false sense of security from your data)
- Our system will have a dashboard where all metrics will be displayed.

## Assumptions
- We're retaining metrics data in our data store for one year.
- The number of hosts sending data is on the upper end of the scale, say 750,000 thousand hosts.
- Each machine inside the company network has an agent installed that streams data to our backend servers.
- The same number of machines will be accessing the metrics dashboard. User will check the dashboard once every hour.

## What are we dealing with?
- Each machine is sending the following (via the agent) to our backend servers:
  - CPU % utilization
  - Memory usage
  - Network I/O
  - Disk I/O
  - GC latency
  - OOMs
  - Number of threads
  - Response time
- Modeling this data as a relational table wouldn't really make sense.
  - Could in theory create a table with a PK of HostID and then rest of the columns, but writing to disk all of these records PER SECOND would be a massive bottleneck
  -  
