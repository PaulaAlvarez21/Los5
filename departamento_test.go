package main

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	db "Los5.com/ServidorWeb/db/sqlc"
)

func datosTestDepartamento(t *testing.T) db.Departamento {
	t.Helper()
	arg := db.CreateDepartamentoParams{
		Nombre:      "Depto Test",
		Direccion:   "Calle Falsa 123",
		Disponible:  true,
		Limpio:      true,
		Descripcion: sql.NullString{String: "Descripcion de prueba", Valid: true},
	}
	depto, err := testQueries.CreateDepartamento(context.Background(), arg)
	if err != nil {
		t.Fatalf("no se pudo crear departamento: %v", err)
	}
	t.Cleanup(func() { testQueries.DeleteDepartamento(context.Background(), depto.IDDepto) })
	return depto
}

func TestCreateDepartamento(t *testing.T) {
	depto := datosTestDepartamento(t)

	if depto.IDDepto == 0 {
		t.Error("id_depto deberia ser generado por la base de datos")
	}
	if depto.Nombre != "Depto Test" {
		t.Errorf("nombre incorrecto, esperado 'Depto Test', got '%s'", depto.Nombre)
	}
	if depto.Direccion != "Calle Falsa 123" {
		t.Errorf("direccion incorrecta, esperado 'Calle Falsa 123', got '%s'", depto.Direccion)
	}
	if !depto.Disponible {
		t.Error("disponible deberia ser true")
	}
	if !depto.Limpio {
		t.Error("limpio deberia ser true")
	}
	if depto.Descripcion.String != "Descripcion de prueba" {
		t.Errorf("descripcion incorrecta, got '%s'", depto.Descripcion.String)
	}
}

func TestGetDepartamento(t *testing.T) {
	depto := datosTestDepartamento(t)

	got, err := testQueries.GetDepartamento(context.Background(), depto.IDDepto)
	if err != nil {
		t.Fatalf("no se pudo obtener departamento: %v", err)
	}
	if got.IDDepto != depto.IDDepto {
		t.Errorf("id_depto incorrecto, esperado %d, got %d", depto.IDDepto, got.IDDepto)
	}
	if got.Nombre != depto.Nombre {
		t.Errorf("nombre incorrecto, esperado '%s', got '%s'", depto.Nombre, got.Nombre)
	}
}

func TestListDepartamentos(t *testing.T) {
	datosTestDepartamento(t)
	datosTestDepartamento(t)

	departamentos, err := testQueries.ListDepartamentos(context.Background())
	if err != nil {
		t.Fatalf("no se pudo listar departamentos: %v", err)
	}
	if len(departamentos) < 2 {
		t.Errorf("esperaba al menos 2 departamentos, got %d", len(departamentos))
	}
}

func TestUpdateDepartamento(t *testing.T) {
	depto := datosTestDepartamento(t)

	arg := db.UpdateDepartamentoParams{
		IDDepto:     depto.IDDepto,
		Nombre:      "Depto Actualizado",
		Direccion:   "Nueva Direccion 456",
		Disponible:  false,
		Limpio:      false,
		Descripcion: sql.NullString{String: "Nueva descripcion", Valid: true},
	}
	err := testQueries.UpdateDepartamento(context.Background(), arg)
	if err != nil {
		t.Fatalf("no se pudo actualizar departamento: %v", err)
	}

	got, err := testQueries.GetDepartamento(context.Background(), depto.IDDepto)
	if err != nil {
		t.Fatalf("no se pudo obtener departamento actualizado: %v", err)
	}
	if got.Nombre != "Depto Actualizado" {
		t.Errorf("nombre no se actualizo, esperado 'Depto Actualizado', got '%s'", got.Nombre)
	}
	if got.Direccion != "Nueva Direccion 456" {
		t.Errorf("direccion no se actualizo, got '%s'", got.Direccion)
	}
	if got.Disponible {
		t.Error("disponible deberia ser false")
	}
}

func TestDeleteDepartamento(t *testing.T) {
	depto := datosTestDepartamento(t)

	err := testQueries.DeleteDepartamento(context.Background(), depto.IDDepto)
	if err != nil {
		t.Fatalf("no se pudo eliminar departamento: %v", err)
	}

	_, err = testQueries.GetDepartamento(context.Background(), depto.IDDepto)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("esperaba sql.ErrNoRows al obtener departamento eliminado, got %v", err)
	}
}
