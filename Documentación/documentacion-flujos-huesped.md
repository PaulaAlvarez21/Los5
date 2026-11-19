# Flujos de creación de huésped

## Alta desde Huéspedes
- Muestra el selector de **Reserva**.
- Flujo: registrar un huésped complementario asociado a una reserva ya existente, eligiendo a cuál corresponde.

## Atajo desde Nueva reserva
- No muestra el selector de **Reserva**.
- Flujo: durante la creación de la reserva, permite crear rápido al huésped principal responsable si se olvidó cargarlo antes; la asociación a la reserva se resuelve en el propio flujo de alta de reserva (selección del huésped), no en este formulario.
- Implementación: el formulario se abre con `?sin_reserva=1` (`hx-get="/huespedes/nuevo?sin_reserva=1"`) y el contexto se preserva en el envío para mantener el modo ante errores. El guardado cierra el modal vía HTMX (`hx-target="#modal-root"`) y, si existe el listado de huéspedes, inserta la fila por OOB.
