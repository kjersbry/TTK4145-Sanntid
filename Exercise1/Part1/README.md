Part1: Thinking about elevators
---------------------------
Not for handing in, just for thinking about. Talk to other groups, assistants, or even people who have taken the course in previous years.

Brainstorm some techniques you could use to prevent a user from being hopelessly stranded, waiting for an elevator that will never arrive. Think about the [worst-case](http://xkcd.com/748/) behaviour of the system.
 - What if the software controlling one of the elevators suddenly crashes?
 > forutsatt at koden for de andre heisene fortsatt kjører: Den som kræsjer kjører sikkerhetsmekanisme (kjører til nærmeste etg og åpner dører), mens de andre tar bestillingen til den heisen
 > Hvis alle kjæsjer: kjøre til nærmeste etg.
 Backup av køen

 - What if it doesn't crash, but hangs?
 > Hvis den kommer seg ut av det den henger i sjekker den om bestillingen er fikset

 - What if a message between machines is lost?
 > For eksempel bruke noe acknowledge greier
 > Heisene har register over hvor mange meldinger de har sendt, sammenligne med et register for alle. Kjøre noe fiksing hvis de ikke stemmer

 - What if the network cable is suddenly disconnected? Then re-connected?
 > Antar: en heis blir koblet fra de andre: Den heisen kjører ferdig bestillingene som er gjort fra innsiden, før den stopper og gir feilsignal.
 > Dersom alle heisene blir koblet fra hverandre og ikke kan kommunisere lenger, switcher de til å fungere som enkeltstående heiser

 - What if a user of the system is being a troll?
 > Systemet kan ikke kræsje selv om alle knapper trykkes osv
 > Håndtere alt på vanlig måte bare?

 - What if the elevator car never arrives at its destination?
> Kan bruke en timer? Never = en bestemt tid. Hvis den ikke er der innen da må den sende bestillingen til en annen heis, og evt gi feilmelding (som kanskje kan trekkes tilbake hvis den oppdager at den funker igjen f eks vha etasjeføler)
