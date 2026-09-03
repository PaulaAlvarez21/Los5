package main

import (
	"context"
	"database/sql"
	"testing"

	db "Los5.com/ServidorWeb/db/sqlc"
)

func createTestDepartamento(t *testing.T) db.Departamento {
	t.Helper()
	arg := db.CreateDepartamentoParams{
		IDDepto:     1,
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
	return depto
}

func TestCreateDepartamento(t *testing.T) {
	depto := createTestDepartamento(t)

	if depto.IDDepto != 1 {
		t.Errorf("id_depto incorrecto, esperado 1, got %d", depto.IDDepto)
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

	// limpiar
	testQueries.DeleteDepartamento(context.Background(), depto.IDDepto)
}

func TestGetDepartamento(t *testing.T) {
	depto := createTestDepartamento(t)

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

	// limpiar
	testQueries.DeleteDepartamento(context.Background(), depto.IDDepto)
}

func TestListDepartamentos(t *testing.T) {
	// crear 2 departamentos
	depto1 := createTestDepartamento(t)
	depto2 := createTestDepartamento(t)

	departamentos, err := testQueries.ListDepartamentos(context.Background())
	if err != nil {
		t.Fatalf("no se pudo listar departamentos: %v", err)
	}
	if len(departamentos) < 2 {
		t.Errorf("esperaba al menos 2 departamentos, got %d", len(departamentos))
	}

	// limpiar
	testQueries.DeleteDepartamento(context.Background(), depto1.IDDepto)
	testQueries.DeleteDepartamento(context.Background(), depto2.IDDepto)
}

func TestUpdateDepartamento(t *testing.T) {
	depto := createTestDepartamento(t)

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

	// limpiar
	testQueries.DeleteDepartamento(context.Background(), depto.IDDepto)
}

func TestDeleteDepartamento(t *testing.T) {
	depto := createTestDepartamento(t)

	err := testQueries.DeleteDepartamento(context.Background(), depto.IDDepto)
	if err != nil {
		t.Fatalf("no se pudo eliminar departamento: %v", err)
	}

	_, err = testQueries.GetDepartamento(context.Background(), depto.IDDepto)
	if err == nil {
		t.Error("esperaba error al obtener departamento eliminado, pero no hubo error")
	}
}
