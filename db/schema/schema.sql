-- Created by Redgate Data Modeler (https://datamodeler.redgate-platform.com)
-- Last modification date: 2026-09-02 22:05:57.984

-- tables
-- Table: departamento
DROP TABLE IF EXISTS huesped;
DROP TABLE IF EXISTS reserva;
DROP TABLE IF EXISTS departamento;

CREATE TABLE departamento (
    id_depto serial  PRIMARY KEY,
    nombre varchar(100)  NOT NULL,
    direccion varchar(150)  NOT NULL,
    disponible boolean  NOT NULL,
    limpio boolean  NOT NULL,
    descripcion varchar(500)  NULL
);

-- Table: reserva
CREATE TABLE reserva (
    id_reserva serial  PRIMARY KEY,
    fecha_inicio date  NOT NULL,
    id_depto int  NOT NULL,
    fecha_fin date  NOT NULL,
    precio_base decimal(10,2)  NOT NULL,
    cant_noches int  NOT NULL,
    descuento decimal(5,2)  NULL,
    observaciones varchar(500)  NULL,
    CONSTRAINT fk_reserva_depto FOREIGN KEY (id_depto)
        REFERENCES departamento (id_depto)
);

-- Table: huesped
CREATE TABLE huesped (
    id_huesped serial  PRIMARY KEY,
    id_reserva int  NOT NULL,
    nombre varchar(50)  NOT NULL,
    apellido varchar(50)  NOT NULL,
    telefono varchar(50)  NULL,
    email varchar(120)  NULL,
    observaciones varchar(255)  NULL,
    CONSTRAINT fk_huesped_reserva FOREIGN KEY (id_reserva)
        REFERENCES reserva (id_reserva)
);

-- End of file.
