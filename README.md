## Project Summary

https://sherrys-shakesearch.onrender.com/

In this project, I aimed to improve the user experience of Shakesearch by providing more context and better UI design.  The first thing I felt was missing from the the original Shakesearch was additional detail on each play. Pulling long excerpts out of context without details on the play or which character was speaking seemed clumsy and difficult to read. I was able to find a csv file from [Kaggle](https://www.kaggle.com/datasets/kingburrito666/shakespeare-plays?resource=download) that organized each line from each of Shakespeare's works by Play, Character, Line, and Act/Scene/Line. I used this file to grab that data from the search term and display it to the user. 

I noticed that the csv file only supplies a single line containing the search term and the original Shakesearch supplied more context around the search term. I didn't want to remove this feature so I created a modal with a "More" button that pulls additional text before and after the search term, using the original text file provided. I also highlighted the search term inside that modal to make it easier for the reader to find what they are looking for.

I also noticed that if a word was not found in the original search, the user is not notifed in any way. I created a modal that lets the user know that no results were found.

Additionally, I used Materialize to make some improvements to the UI design. One major change was switching from a table to cards because I felt it was more interesting to read through that way.

# Further Improvements

One issue I noticed was that there are slight differences in the way some lines are transcribed in the csv file versus the text file. For instance, a comma might be missing or a word is spelled differently. This prevents the search engine from being able to provide additional context when the "More" button is pushed for certain search terms. With more time, I would improve the data files to make sure they matched precisely and context could be found for every quote.

I would add suggestions for misspelled words in the search engine.

I would work with a designer to further improve the UI design.

This was my first time working in Go so I would ask for review from someone more experienced to help me with following best practices.

There are areas where I could improve performance, for instance using an improved search algorithm from the linear `suffixarray` search. I could also condiser using a cache rather than loading the entire works on each search. I could also improve memory allocation by reusing the same buffer and encoder across requests.


# ShakeSearch

Welcome to the Pulley Shakesearch Take-home Challenge! In this repository,
you'll find a simple web app that allows a user to search for a text string in
the complete works of Shakespeare.

You can see a live version of the app at
https://pulley-shakesearch.onrender.com/. Try searching for "Hamlet" to display
a set of results.

In it's current state, however, the app is in rough shape. The search is
case sensitive, the results are difficult to read, and the search is limited to
exact matches.

## Your Mission

Improve the app! Think about the problem from the **user's perspective**
and prioritize your changes according to what you think is most useful.

You can approach this with a back-end, front-end, or full-stack focus.

## Evaluation

We will be primarily evaluating based on how well the search works for users. A search result with a lot of features (i.e. multi-words and mis-spellings handled), but with results that are hard to read would not be a strong submission.

## Submission

1. Fork this repository and send us a link to your fork after pushing your changes.
2. Render (render.com) hosting, the application deploys cleanly from a public url.
3. In your submission, share with us what changes you made and how you would prioritize changes if you had more time.
