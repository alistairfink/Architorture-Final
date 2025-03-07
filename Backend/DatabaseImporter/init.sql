CREATE DATABASE Architorture;
\c architorture

CREATE TABLE public.expansion
(
    id integer NOT NULL,
    name text COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT expansion_pkey PRIMARY KEY (id)
) TABLESPACE pg_default;

ALTER TABLE public.expansion
    OWNER to postgres;

CREATE TABLE public.card_type
(
    id integer NOT NULL,
    name text COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT card_type_pkey PRIMARY KEY (id)
) TABLESPACE pg_default;

ALTER TABLE public.card_type
    OWNER to postgres;

CREATE TABLE public.card
(
    id integer NOT NULL,
    card_type_id integer NOT NULL,
    name text COLLATE pg_catalog."default" NOT NULL,
    description text COLLATE pg_catalog."default",
    play_immediately boolean NOT NULL,
    quantity integer NOT NULL,
    expansion_id integer NOT NULL,
    archivable boolean NOT NULL,
    CONSTRAINT card_pkey PRIMARY KEY (id),
    CONSTRAINT card_type_id_fk FOREIGN KEY (card_type_id)
        REFERENCES public.card_type (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT expansion_id_fk FOREIGN KEY (expansion_id)
        REFERENCES public.expansion (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

    TABLESPACE pg_default;

ALTER TABLE public.card
    OWNER to postgres;


INSERT INTO public.card_type ("id", "name") VALUES (1, 'Architorture'), (2, 'Save'), (3, 'Primary'), (4, 'Secondary'), (5, 'Tertiary'), (6, 'TertiaryWildCard');

INSERT INTO public.expansion ("id", "name") VALUES (1, 'Base Cards - Year 1 '), (2, 'Expansion 1 - Year 2 '), (3, 'Expansion 2 - Year 3 '), (4, 'Expansion 3 - Year 4');

INSERT INTO public.card (
    "id", "card_type_id", "name", "description",
    "play_immediately", "quantity",
    "expansion_id", "archivable"
)
VALUES
    (
        1, 2, 'Eureka', 'With these cards you are able to dismiss an Architorture card and avoid being eliminated from the game. These are a priceless asset in the game of Architorture.',
        FALSE, 6, 1, TRUE
    ),
    (
        2, 2, 'Eureka', 'With these cards you are able to dismiss an Architorture card and avoid being eliminated from the game. These are a priceless asset in the game of Architorture.',
        FALSE, 1, 4, TRUE
    ),
    (
        3, 1, 'Architorture', 'If you draw these cards, you''ve been overwhelmed by Architorture! This the end of the road for your Architecture education, you are eliminated.',
        FALSE, 4, 1, TRUE
    ),
    (
        4, 1, 'Architorture', 'If you draw these cards, you''ve been overwhelmed by Architorture! This the end of the road for your Architecture education, you are eliminated.',
        FALSE, 1, 3, TRUE
    ),
    (
        5, 3, 'Architecture Dump', 'Discard 1 card that you won’t believe will be useful in the game. This card must be played immediately.',
        TRUE, 6, 1, TRUE
    ),
    (
        6, 3, 'Architecture Dump', 'Discard 1 card that you won’t believe will be useful in the game. This card must be played immediately.',
        TRUE, 2, 4, TRUE
    ),
    (
        7, 3, 'Architecture Memory Loss',
        'You never know when memory loss will strike! The person on your right will blindly choose 1 card from your hand and insert it back into the draw pile. This card must be played immediately.',
        TRUE, 4, 1, TRUE
    ),
    (
        8, 3, 'Architecture Memory Loss',
        'You never know when memory loss will strike! The person on your right will blindly choose 1 card from your hand and insert it back into the draw pile. This card must be played immediately.',
        TRUE, 2, 4, TRUE
    ),
    (
        9, 3, 'Resource Shuffle', 'Everyone in the game must combine the cards in hand together, shuffle and redistribute amongst the players evenly. Deal the cards in a clockwise rotation starting with the player that drew this card.',
        FALSE, 4, 1, TRUE
    ),
    (
        10, 3, 'Resource Shuffle', 'Everyone in the game must combine the cards in hand together, shuffle and redistribute amongst the players evenly. Deal the cards in a clockwise rotation starting with the player that drew this card. ',
        FALSE, 2, 4, TRUE
    ),
    (
        11, 3, 'Memory Expansion Card', 'With this card you are able to +1 card to your hand as your memory has expanded.',
        FALSE, 2, 1, TRUE
    ),
    (
        12, 3, 'Memory Expansion Card', 'With this card you are able to +1 card to your hand as your memory has expanded.',
        FALSE, 1, 3, TRUE
    ),
    (
        13, 3, 'Memory Expansion Card', 'With this card you are able to +1 card to your hand as your memory has expanded.',
        FALSE, 1, 4, TRUE
    ),
    (
        14, 3, 'USB Expansion Card', 'With this card you are able to +1 card to your USB as your USB memory has been expanded. Place this card with the archived card(s) on the table.',
        FALSE, 2, 1, TRUE
    ),
    (
        15, 3, 'USB Expansion Card', 'With this card you are able to +1 card to your USB as your USB memory has been expanded. Placed this card with the archived card(s) on the table.',
        FALSE, 1, 2, TRUE
    ),
    (
        16, 3, 'USB Expansion Card', 'With this card you are able to +1 card to your USB as your USB memory has been expanded. Placed this card with the archived card(s) on the table.',
        FALSE, 1, 4, TRUE
    ),
    (
        17, 3, 'Found a USB', 'Congratulations! You found a forgotten USB in the computer lab. Take an USB Expansion Card from another player. If no expansion cards are played, keep card and play at the next available opportunity. ',
        FALSE, 1, 2, TRUE
    ),
    (
        18, 3, 'Found a USB', 'Congratulations! You found a forgotten USB in the computer lab. Take an USB Expansion Card from another player. If no expansion cards are played, keep card and play at the next available opportunity. ',
        FALSE, 1, 4, TRUE
    ),
    (
        19, 4, 'Vendetta', 'Sometimes things don’t go your way but it doesn’t mean your classmates shouldn’t have one. Give the next player a card from your hand and force them to draw 2 additional cards to end their turn.',
        FALSE, 4, 1, TRUE
    ),
    (
        20, 4, 'Preview', 'You have an opportunity to prepare and take a glimpse of the future. Draw the first two (2) cards and preview what the next resources or downfall cards will be. Place the cards back on top of the draw pile in the same order once you’re done viewing them.',
        FALSE, 5, 1, TRUE
    ),
    (
        21, 4, 'Extra Preview', 'You have an opportunity to EXTRA prepare and take a glimpse of the future. Draw the first four (4) cards and preview what the next resources or downfall cards will be. Place the cards back on top of the draw pile in the same order once you’re done viewing them.',
        FALSE, 3, 3, TRUE
    ),
    (
        22, 4, 'Thank U, Next!', 'Want to avoid the next draw? Play this card and skip your turn. Thank U, Next!',
        FALSE, 4, 1, TRUE
    ),
    (
        23, 4, 'Shuffle', 'Shuffle the cards in the draw file.',
        FALSE, 4, 1, TRUE
    ),
    (
        24, 4, 'Reverse', 'Not liking the order of the players? Let''s switch it up. Use this card to reverse the order and end your turn.',
        FALSE, 2, 2, TRUE
    ),
    (
        25, 4, 'Change of Plans', 'Sometimes you got to look at your different options and get your ducks in a row. You have an opportunity to change the near future. Take a peek at the first two (2) top cards, rearrange in the order you wish them to be and place them back onto the draw pile once you''re done. ',
        FALSE, 3, 2, TRUE
    ),
    (
        26, 4, 'Change of Plans', 'Sometimes you got to look at your different options and get your ducks in a row. You have an opportunity to change the near future. Take a peek at the first two (2) top cards, rearrange in the order you wish them to be and place them back onto the draw pile once you''re done. ',
        FALSE, 3, 4, TRUE
    ),
    (
        27, 4, 'Not Today', 'Did a student or professor try to pass a fast one on you? Well today is your lucky day! Play this card to cancel out another player''s action card. The card must be played before the action is followed through and could be played at anytime (even when it’s not your turn). Unfortunately, Architorture cards cannot be voided. ',
        FALSE, 5, 2, TRUE
    ),
    (
        28, 4, 'Not Today', 'Did a student or professor try to pass a fast one on you? Well today is your lucky day! Play this card to cancel out another player''s action card(s). The card must be played before the action is followed through and could be played at anytime (even when it’s not your turn). Unfortunately, Architorture cards cannot be voided. ',
        FALSE, 4, 4, TRUE
    ),
    (
        29, 6, 'Surprise Epiphany', 'Sometimes good ideas can’t wait. Use this card as a substitute of one of your knowledge cards to complete a pairing or trioing play.',
        FALSE, 2, 2, TRUE
    ),
    (
        30, 6, 'Surprise Epiphany', 'Sometimes good ideas can’t wait. Use this card as a substitute of one of your knowledge cards to complete a pairing or trioing play.',
        FALSE, 2, 3, TRUE
    ),
    (
        31, 4, 'Super Vendetta', 'This card is similar to the Vendetta except you get to choose who you want to attack. Give the player of your choice a card from your hand and force them to draw 2 additional cards to end their turn.',
        FALSE, 3, 2, TRUE
    ),
    (
        32, 4, 'Super Vendetta', 'This card is similar to the Vendetta except you get to choose who you want to attack. Give the player of your choice a card from your hand and force them to draw 2 additional cards to end their turn.',
        FALSE, 1, 3, TRUE
    ),
    (
        33, 4, 'Co-op', 'Time to take a break. Play this card to skip the next two rounds of your turn. ',
        FALSE, 2, 3, TRUE
    ),
    (
        34, 4, 'Not in this Lifetime', 'Nope their nope, the nope of ALL nopes! The Lord has spoken, god says NO! :) Play this card to cancel out another player''s action card(s) and overpower any Not Today card. The card must be played before the action is followed through and could be played at anytime (even when it’s not your turn). Unfortunately, Architorture cards cannot be voided. ',
        FALSE, 1, 3, TRUE
    ),
    (
        35, 2, 'Time Extension', 'Glory be to the Architecture Gods, your professor has decided to be kind to you and give you an extension. You are able to avoid the wrath of architorture. In the future, if an Architorture card is drawn, you may hold the Architorture card in your hand as long as your have the Time Extension card there as well. This card cannot be archived and must be in your hand at all times and take up a spot in your memory.',
        FALSE, 2, 3, FALSE
    ),
    (
        36, 4, 'Share the Preview', 'Draw the top three (3) cards and share the preview with another player of your choice. Place the cards back on top of the draw pile in the same order once you’re done viewing them.',
        FALSE, 3, 4, TRUE
    ),
    (
        37, 2, 'Numb Card', 'You’ve gotten to the point where you''re just numb, you stopped feeling the architorture. Use this card to deflect the architorture card toward another player of your choosing.',
        FALSE, 1, 4, TRUE
    ),
    (
        38, 4, 'Assistance', 'Play this card to ask for assistance from another player of your choosing. That player must give you one of their cards of choice.',
        FALSE, 4, 4, TRUE
    ),
    (
        39, 5, 'Curriculum', 'Knowledge cards will allow you to steal resources from other students throughout the game. Pairs or trios of the same cards will allow you to make power moves in the game. ',
        FALSE, 5, 1, TRUE
    ),
    (
        40, 5, 'Curriculum', 'Knowledge cards will allow you to steal resources from other students throughout the game. Pairs or trios of the same cards will allow you to make power moves in the game. ',
        FALSE, 1, 2, TRUE
    ),
    (
        41, 5, 'Curriculum', 'Knowledge cards will allow you to steal resources from other students throughout the game. Pairs or trios of the same cards will allow you to make power moves in the game. ',
        FALSE, 1, 3, TRUE
    ),
    (
        42, 5, 'Curriculum', 'Knowledge cards will allow you to steal resources from other students throughout the game. Pairs or trios of the same cards will allow you to make power moves in the game. ',
        FALSE, 2, 4, TRUE
    ),
    (
        43, 5, 'Building Knowledge', 'Knowledge cards will allow you to steal resources from other students throughout the game. Pairs or trios of the same cards will allow you to make power moves in the game. ',
        FALSE, 4, 2, TRUE
    ),
    (
        44, 5, 'Building Knowledge', 'Knowledge cards will allow you to steal resources from other students throughout the game. Pairs or trios of the same cards will allow you to make power moves in the game. ',
        FALSE, 2, 3, TRUE
    ),
    (
        45, 5, 'Building Knowledge', 'Knowledge cards will allow you to steal resources from other students throughout the game. Pairs or trios of the same cards will allow you to make power moves in the game. ',
        FALSE, 2, 4, TRUE
    ),
    (
        46, 5, 'Software Skills', 'Knowledge cards will allow you to steal resources from other students throughout the game. Pairs or trios of the same cards will allow you to make power moves in the game. ',
        FALSE, 3, 2, TRUE
    ),
    (
        47, 5, 'Software Skills', 'Knowledge cards will allow you to steal resources from other students throughout the game. Pairs or trios of the same cards will allow you to make power moves in the game. ',
        FALSE, 3, 3, TRUE
    ),
    (
        48, 5, 'Software Skills', 'Knowledge cards will allow you to steal resources from other students throughout the game. Pairs or trios of the same cards will allow you to make power moves in the game. ',
        FALSE, 1, 4, TRUE
    ),
    (
        49, 5, 'Time Extension', 'Knowledge cards will allow you to steal resources from other students throughout the game. Pairs or trios of the same cards will allow you to make power moves in the game. ',
        FALSE, 5, 1, TRUE
    ),
    (
        50, 5, 'Starchitects Knowledge', 'Knowledge cards will allow you to steal resources from other students throughout the game. Pairs or trios of the same cards will allow you to make power moves in the game. ',
        FALSE, 2, 3, TRUE
    ),
    (
        51, 5, 'Starchitects Knowledge', 'Knowledge cards will allow you to steal resources from other students throughout the game. Pairs or trios of the same cards will allow you to make power moves in the game. ',
        FALSE, 2, 4, TRUE
    ),
    (
        52, 5, 'Architectural Motor Skills',
        'Knowledge cards will allow you to steal resources from other students throughout the game. Pairs or trios of the same cards will allow you to make power moves in the game. ',
        FALSE, 5, 1, TRUE
    ),
    (
        53, 5, 'Architectural Motor Skills',
        'Knowledge cards will allow you to steal resources from other students throughout the game. Pairs or trios of the same cards will allow you to make power moves in the game. ',
        FALSE, 2, 3, TRUE
    ),
    (
        54, 5, 'Architectural Motor Skills',
        'Knowledge cards will allow you to steal resources from other students throughout the game. Pairs or trios of the same cards will allow you to make power moves in the game. ',
        FALSE, 1, 4, TRUE
    ),
    (
        55, 5, 'Personal Growth', 'Knowledge cards will allow you to steal resources from other students throughout the game. Pairs or trios of the same cards will allow you to make power moves in the game. ',
        FALSE, 5, 1, TRUE
    ),
    (
        56, 5, 'Personal Growth', 'Knowledge cards will allow you to steal resources from other students throughout the game. Pairs or trios of the same cards will allow you to make power moves in the game. ',
        FALSE, 2, 3, TRUE
    ),
    (
        57, 5, 'Personal Growth', 'Knowledge cards will allow you to steal resources from other students throughout the game. Pairs or trios of the same cards will allow you to make power moves in the game. ',
        FALSE, 2, 4, TRUE
    );
