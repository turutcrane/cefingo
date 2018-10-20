#include "cefingo.h"
#include "_cgo_export.h"

void construct_cefingo_life_span_handler(cefingo_life_span_handler_wrapper_t *handler) {
    // cefingo_cslogf(__func__, "L14: 0x%llx", handler);
    initialize_cefingo_base_ref_counted(
        offsetof(cefingo_life_span_handler_wrapper_t, counter),
        (cef_base_ref_counted_t*) handler);
    handler->body.on_before_close = life_span_on_before_close;
    handler->body.do_close = life_span_do_close;
    handler->body.on_after_created = life_span_on_after_created;
}

void construct_cefingo_browser_process_handler(cefingo_browser_process_handler_wrapper_t *handler) {
    initialize_cefingo_base_ref_counted(
        offsetof(cefingo_browser_process_handler_wrapper_t, counter),
        (cef_base_ref_counted_t*) handler);
    handler->body.on_context_initialized = browser_process_on_context_initialized;
}

void construct_cefingo_client(cefingo_client_wrapper_t* client) {
    initialize_cefingo_base_ref_counted(
        offsetof(cefingo_client_wrapper_t, counter),
        (cef_base_ref_counted_t*) client);

//     // callbacks
    client->body.get_context_menu_handler = get_context_menu_handler;
    client->body.get_dialog_handler = get_dialog_handler;
    client->body.get_display_handler = get_display_handler;
    client->body.get_download_handler = get_download_handler;
    client->body.get_drag_handler = get_drag_handler;
    client->body.get_find_handler = get_find_handler;
    client->body.get_focus_handler = get_focus_handler;
    client->body.get_jsdialog_handler = get_jsdialog_handler;
    client->body.get_keyboard_handler = get_keyboard_handler;
    client->body.get_life_span_handler = get_life_span_handler;  // Implemented!
    client->body.get_load_handler = client_get_load_handler;
    client->body.get_render_handler = get_render_handler;
    client->body.get_request_handler = get_request_handler;
    client->body.on_process_message_received = client_on_process_message_received;
}

static void CEF_CALLBACK c_on_before_command_line_processing(
        struct _cef_app_t* self, const cef_string_t* process_type,
        struct _cef_command_line_t* command_line) {

    // simply call go func with type cast
    on_before_command_line_processing(self, (cef_string_t *) process_type, command_line);
}

void construct_cefingo_app(cefingo_app_wrapper_t* app) {
    initialize_cefingo_base_ref_counted(
        offsetof(cefingo_app_wrapper_t, counter),
        (cef_base_ref_counted_t*) app);

    // callbacks
    app->body.on_before_command_line_processing = c_on_before_command_line_processing;
    app->body.on_register_custom_schemes = on_register_custom_schemes;

    app->body.get_resource_bundle_handler = get_resource_bundle_handler;
    app->body.get_browser_process_handler = get_browser_process_handler;
    app->body.get_render_process_handler = get_render_process_handler;
}

void construct_cefingo_render_process_handler(cefingo_render_process_handler_wrapper_t* handler) {
    initialize_cefingo_base_ref_counted(
        offsetof(cefingo_render_process_handler_wrapper_t, counter),
        (cef_base_ref_counted_t*) handler);

    // callbacks
    handler->body.on_render_thread_created = on_render_thread_created;
    handler->body.on_context_created = on_context_created;
    handler->body.on_web_kit_initialized = on_web_kit_initialized;
    handler->body.on_browser_created = on_browser_created;
    handler->body.on_browser_destroyed = on_browser_destroyed;
    handler->body.get_load_handler = render_process_hander_get_load_handler;
    handler->body.on_context_released = on_context_released;
    handler->body.on_uncaught_exception = on_uncaught_exception;
    handler->body.on_focused_node_changed = on_focused_node_changed;
    handler->body.on_process_message_received = render_process_handler_on_process_message_received;

}