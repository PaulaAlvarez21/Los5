# Flujos de creación de huésped

El modal de alta de huésped (`views/modal.templ` → `HuespedFormModal` →
`HuespedForm`) se usa en dos contextos. La diferencia se decide en la apertura
con el parámetro `?sin_reserva=1` y se conserva en el envío mediante el campo
oculto `sin_reserva`, que **siempre** está presente en el formulario:

| Origen | `sin_reserva` | Selector de **Reserva** |
|---|---|---|
| Botón "+ Nuevo huésped" de `/huespedes` | `false` | Sí |
| Botón "+ Crear huésped" del formulario de reserva | `true` | No |

`NuevoHuespedWeb` usa el parámetro de la URL para elegir el modo, y
`renderHuespedModalError` lo conserva ante errores leyendo tanto la query como
el formulario, para que el modal vuelva a mostrarse igual.

## Alta desde Huéspedes

- Muestra el selector de **Reserva**.
- Flujo: registrar un huésped complementario asociado a una reserva ya existente, eligiendo a cuál corresponde.
- Guardado: **no redirecciona**. Con HTMX, `CrearHuespedWeb` responde
  `views.HuespedCreado`, que inserta la fila nueva en `#huesped-rows` por OOB y
  deja el modal cerrado. Sin HTMX, responde `303` a `/huespedes`.

## Atajo desde Nueva reserva

- No muestra el selector de **Reserva**.
- Flujo: durante la creación de la reserva, permite crear rápido al huésped principal responsable si se olvidó cargarlo antes; la asociación a la reserva se resuelve en el propio flujo de alta de reserva (selección del huésped), no en este formulario.
- El formulario se abre con `?sin_reserva=1`
  (`hx-get="/huespedes/nuevo?sin_reserva=1"`, en `views/entity_form.templ`).
- Guardado: después de crear el huésped, `CrearHuespedWeb` lee `sin_reserva`
  con `parseBoolForm` y, si es `true` y la request es HTMX, responde
  `views.ReabrirFormReserva()`. Ese componente es solo un disparador:

  ```html
  <div hx-get="/reservas/nuevo" hx-trigger="load"
       hx-target="#modal-root" hx-swap="innerHTML"></div>
  ```

  Como el form de huésped ya apunta a `#modal-root`, el disparador reemplaza al
  modal; al cargarse pide `GET /reservas/nuevo`, que devuelve el formulario de
  reserva con datos frescos. El huésped recién creado aparece ya en el combo
  **Huésped principal**, listo para seleccionarlo y continuar. Es un request
  extra, pero evita duplicar la lógica de armado del modal fuera de
  `ReservaWebHandler` y no requiere inyectar servicios nuevos.
- Sin HTMX, el guardado responde `303` a `/reservas` (no se puede reabrir un
  modal sin JS; el usuario debe hacer clic en "+ Nueva reserva").
- El alta rápida **no** inserta la fila en el padrón de huéspedes: esa tabla no
  está visible en la página de reservas. El huésped queda registrado y se ve al
  volver a `/huespedes`.
