drop database if exists payment_db;
create database payment_db;

drop database if exists vehicles_reservations_db;
create database vehicles_reservations_db;

drop database if exists common_db;
create database common_db;


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


--------------------------------------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------------------------------------



--- DATA
USE users_db;

-- Insert fake data for users
INSERT INTO users (Name, Email, contactNo, hashedPassword, membershipTier) VALUES
('John Doe', 'john.doe@example.com', '12345678', 'hashedpassword123', 'Basic'),
('Jane Smith', 'jane.smith@example.com', '87654321', 'hashedpassword456', 'Premium'),
('Alice Johnson', 'alice.johnson@example.com', '23456789', 'hashedpassword789', 'VIP'),
('Bob Brown', 'bob.brown@example.com', '34567890', 'hashedpassword012', 'Basic');


--------------------------------------------------------------------------------------------------------------------


USE common_db;

-- Insert fake data for member benefits
INSERT INTO member_benefits (membershipTier, Name, Description) VALUES
('Basic', 'Free Vehicle Access', 'Access to basic tier vehicles.'),
('Premium', 'Discounted Rates', '10% discount on all rentals.'),
('VIP', 'Priority Support', '24/7 priority customer support.');

-- Insert fake data for promotions
INSERT INTO promotions (Name, Description, Discount, ifPercentage) VALUES
('New Year Promo', 'Get 20% off on your next booking.', 20.00, '1'),
('Flat Discount', 'Get $10 off on your next booking.', 10.00, '0'),
('Holiday Special', '15% off on bookings during the holiday season.', 15.00, '1');


--------------------------------------------------------------------------------------------------------------------


USE vehicles_reservations_db;

-- Insert fake data for vehicles
INSERT INTO vehicles (licensePlate, Model, rentalRate) VALUES
('ABC1234', 'Tesla Model 3', 50),
('DEF5678', 'Nissan Leaf', 40),
('GHI9012', 'Chevy Bolt', 45),
('JKL3456', 'BMW i3', 60);

-- Insert fake data for vehicle status history
INSERT INTO vehicleStatusHistory (vehicleID, location, chargeLevel, cleanlinessStatus) VALUES
(1, 'Downtown Garage', 85, 'Clean'),
(2, 'Airport Lot', 60, 'Clean'),
(3, 'Mall Parking', 20, 'Dirty'),
(4, 'Suburban Garage', 95, 'Clean');

-- Insert fake data for bookings
INSERT INTO bookings (vehicleID, userID, startTime, endTime, Status) VALUES
(1, 1, '2024-12-01 10:00:00', '2024-12-01 12:00:00', 'Completed'),
(2, 2, '2024-12-02 14:00:00', '2024-12-02 16:00:00', 'Active'),
(3, 3, '2024-12-03 09:00:00', '2024-12-03 11:00:00', 'Cancelled'),
(4, 4, '2024-12-04 15:00:00', '2024-12-04 17:00:00', 'Active');


--------------------------------------------------------------------------------------------------------------------

USE payment_db;

-- Insert fake data for payments
INSERT INTO payments (userID, bookingID, Status, promotionID, Amount) VALUES
(1, 1, 'Successful', 1, 40.00),
(2, 2, 'Pending', 2, 30.00),
(3, 3, 'Refunded', 3, 38.25),
(4, 4, 'Successful', 1, 48.00);
