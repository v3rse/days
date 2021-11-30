# A Counter of Days

```txt
 9 For all our days pass away under your wrath;
    we bring our years to an end like a sigh.
10 The years of our life are seventy,
    or even by reason of strength eighty;
yet their span[a] is but toil and trouble;
    they are soon gone, and we fly away.
11 Who considers the power of your anger,
    and your wrath according to the fear of you?
12 So teach us to number our days
    that we may get a heart of wisdom.

                                    -Psalm 90:9-12 (ESV)
```



## Assumptions
Humans live averagely for __70 years__

## Features
- [x] `days journal write`
    - create journal entry
- [ ] `days journal read [start date] [end date]`
    - _no args_: show entries for the day by the hour
    - _start date_: show all entries after this date by the hour
    - _start date, end date_: show all entries between these date by the hour
- [x] `days life start "<birth-date>"` __NEW__
    - sets your birth date
- [x] `days life end [-v]` __NEW__
    - prints summary of days spent alive and estimated days remaining (add the `-v` flag to see some more details)
- [x] `days track <activity>`
    - tracks an activity/habit
- [x] `days since "<activity>"`
    - prints days since tracking started
- [x] `days list`
    - list all activities and day since tracking started
- [x] `days reset "<activity>"`
    - reset tracker for given activity

