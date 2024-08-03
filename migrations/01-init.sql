create table migrations (
    name text not null,
    timestamp int not null,
    primary key (name)
);

create table candlestick_data (
    timestamp int not null,

    bid_o real not null,
    bid_h real not null,
    bid_l real not null,
    bid_c real not null,

    mid_o real not null,
    mid_h real not null,
    mid_l real not null,
    mid_c real not null,

    ask_o real not null,
    ask_h real not null,
    ask_l real not null,
    ask_c real not null,

    primary key (timestamp)
);
