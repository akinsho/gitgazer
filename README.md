### Gitgazer

**This project is archived, I didn't realise when I started it that there is already a much more established and functional project which does most of this
ie. [`gh-dash`](https://github.com/dlvhdr/gh-dash) as well as some equally impressive up and comers, [`gh-eco`](https://github.com/coloradocolby/gh-eco)**

<img width="1720" alt="Screen Shot 2022-04-14 at 14 56 47" src="https://user-images.githubusercontent.com/22454918/163406302-018e9e60-8fab-4c72-ab95-089b6063a532.png">

Gitgazer is a TUI application that allows you keep track of the repos you ~~spy~~ 👉 _gaze_ upon the most.
You can see how the star count has changed over the course of a day or week or month (TBD).
You can see the most recent issues as well as the issues you've tagged/tracked.

## Status

This project is currently in **early development**. It is not ready for daily use.

## Goals

- [x] Decide on main layout for the application
- [x] Handle authentication
- [x] list repositories you've starred (paginated)
- [x] Select a subset of repositories as your focused ones
- [x] Store this list locally so it can be retrieved later
- [x] Persist credentials
- [x] Create config file
- [ ] Make repo list sort order consistent
- [ ] Track the star count for a repository by time window e.g. day, month, year
- [ ] Visualise star count graphically
- [ ] See issues you are watching and track updates since you last checked
