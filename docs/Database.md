# Database
# user-service

## User
```
_id: ObjectID
name: string
email: string
password: string
betaUser: bool
created_at_ts: datetime
updated_at_ts: datetime
deleted_at_ts: datetime
```

### indices:
```
{
    key: email,
    order: ascending
    unique: true
}
```

# deck-managament-service

## Deck
```
_id: ObjectID
name: string
description: string
color: string
user_id: string
cards: []card
created_at: datetime
updated_at: datetime
deleted_at: datetime
```

## Card
```
_id: ObjectID
question: string
answer: string
user_id: string
deck_id: string
created_at: datetime
updated_at: datetime
deleted_at: datetime
```

### indices

```
{
    key: user_id
    order: ascending
}
```

# config-service

## Deck Config
```
_id: ObjectID
name: string
colors: []string
```

# learning-service

## LearningSession
```
_id: ObjectID
user_id: string
deck_id: string
started_at: datetime
finished_at: datetime
finished: bool
```

## CardEvent
```
_id: ObjectID
user_id: string
deck_id: string
card_id: string
learning_session_id: string
memory_half_time: float
number_practiced: int
number_correct: int
number_incorrect: int
number_practiced_last_session: int
number_correct_last_session: int
number_incorrect_last_session: int
created_at: datetime
started_at: datetime
finished_at: datetime
```

### indices
```
{
    keys: user_id, card_id
    order: ascending
}
{
    keys: user_id, deck_id
    order: ascending
}
{
    key: created_at
    order: ascending
}
```

# card-generation-service
WIP

## Note

```
_id: ObjectID
encoding: string
user_id: string
deck_id: string
text: string
completion: string
cards_added: boolean
cards: []Card
creaed_at: datetime
```

## Card
```
cards: []{
    question: string
    answer: string
}
cards_added: boolean
original_cards: boolean
created_at: datetime
```
