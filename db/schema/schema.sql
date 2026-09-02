-- Created by Redgate Data Modeler (https://datamodeler.redgate-platform.com)
-- Last modification date: 2026-09-02 22:05:57.984

-- tables
-- Table: DEPARTAMENTO
CREATE TABLE DEPARTAMENTO (
    id_depto int  NOT NULL,
    nombre varchar(100)  NOT NULL,
    direccion varchar(150)  NOT NULL,
    disponible boolean  NOT NULL,
    limpio boolean  NOT NULL,
    descripcion varchar(500)  NULL,
    CONSTRAINT PK_DEPARTAMENTO PRIMARY KEY (id_depto)
);

-- Table: HUESPED
CREATE TABLE HUESPED (
    id_huesped int  NOT NULL,
    id_reserva int  NOT NULL,
    nombre varchar(50)  NOT NULL,
    apellido varchar(50)  NOT NULL,
    telefono varchar(50)  NULL,
    email varchar(120)  NULL,
    observaciones varchar(255)  NULL,
    CONSTRAINT PK_HUESPED PRIMARY KEY (id_huesped)
);

-- Table: RESERVA
CREATE TABLE RESERVA (
    id_reserva int  NOT NULL,
    fecha_inicio date  NOT NULL,
    id_depto int  NOT NULL,
    fecha_fin date  NOT NULL,
    precio_base decimal(10,2)  NOT NULL,
    cant_noches int  NOT NULL,
    descuento decimal(5,2)  NULL,
    observaciones varchar(500)  NULL,
    CONSTRAINT PK_RESERVA PRIMARY KEY (id_reserva)
);

-- foreign keys
-- Reference: FK_HUESPED_RESERVA (table: HUESPED)
ALTER TABLE HUESPED ADD CONSTRAINT FK_HUESPED_RESERVA
    FOREIGN KEY (id_reserva)
    REFERENCES RESERVA (id_reserva)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: FK_RESERVA_DEPARTAMENTO (table: RESERVA)
ALTER TABLE RESERVA ADD CONSTRAINT FK_RESERVA_DEPARTAMENTO
    FOREIGN KEY (id_depto)
    REFERENCES DEPARTAMENTO (id_depto)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- End of file.

