drop database if exists users_db;
create database users_db;

use users_db;

create table users (
	userID int auto_increment primary key,
    Name varchar(100) not null,
    Email varchar(100) unique not null,
    contactNo char(8) not null,
    hashedPassword varchar(255) not null,
    membershipTier enum('Basic','Premium','VIP') default 'Basic'
);


--------------------------------------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------------------------------------


drop database if exists common_db;
create database common_db;

use common_db;

create table member_benefits (
	benefitID int auto_increment primary key,
    membershipTier enum('Basic', 'Premium', 'VIP') not null,
    Name varchar(255) not null,
    Description TEXT not null
);

create table promotions (
	promotionID int auto_increment primary key,
    Name varchar(100) not null,
    Description varchar(255) not null,
    Discount decimal(10,2) not null,
    ifPercentage enum ('1','0') not null
);


--------------------------------------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------------------------------------


drop database if exists vehicles_reservations_db;
create database vehicles_reservations_db;

use vehicles_reservations_db;

create table vehicles (
	vehicleID int auto_increment primary key,
    licensePlate varchar(20) unique not null,
    Model varchar(100) not null,
    rentalRate int not null
);

create table vehicleStatusHistory (
	statusID int auto_increment primary key,
    vehicleID int not null,
    timestamp TIMESTAMP default CURRENT_TIMESTAMP,
    location varchar(255) not null,
    chargeLevel int not null,
    cleanlinessStatus enum ('Clean','Dirty') not null,
    foreign key (vehicleID) references vehicles(vehicleID)
);

create table bookings (
	bookingID int auto_increment primary key,
    vehicleID int not null,
    userID int not null,
    startTime datetime not null,
    endTime datetime not null,
    Status enum('Active','Completed','Cancelled') default 'Active',
    createdAt timestamp default current_timestamp,
    foreign key (userID) references users_db.users(userID),
    foreign key (vehicleID) references vehicles(vehicleID)
);


--------------------------------------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------------------------------------


drop database if exists payment_db;
create database payment_db;

use payment_db;

create table payments (
	paymentID int auto_increment primary key,
    userID int not null,
    bookingID int not null,
    Status enum('Pending','Successful','Refunded','Unsuccessful') not null default 'Pending',
    promotionID int not null,
    Amount decimal(10,2) not null,
    createdAt timestamp default current_timestamp,
    foreign key (userID) references users_db.users(userID),
    foreign key (bookingID) references vehicles_reservations_db.bookings(bookingID),
    foreign key (promotionID) references common_db.promotions(promotionID)
);