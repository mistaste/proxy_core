cmake_minimum_required(VERSION 3.14)

set(PROJECT_NAME "proxy_core")
project(${PROJECT_NAME} LANGUAGES CXX)

# Name used for the Flutter plugin shared library.
set(PLUGIN_NAME "proxy_core_plugin")

# Unzip libproxy_core.zip (shipped Go shared library) once per configure.
set(PROXY_CORE_ZIP "${CMAKE_CURRENT_SOURCE_DIR}/libproxy_core.zip")
set(PROXY_CORE_DLL "${CMAKE_CURRENT_SOURCE_DIR}/libproxy_core.dll")

if(NOT EXISTS ${PROXY_CORE_ZIP})
    message(FATAL_ERROR "libproxy_core.zip not found at ${PROXY_CORE_ZIP}")
endif()

message(STATUS "Extracting libproxy_core.zip...")
execute_process(
    COMMAND ${CMAKE_COMMAND} -E tar xvf "${PROXY_CORE_ZIP}"
    WORKING_DIRECTORY "${CMAKE_CURRENT_SOURCE_DIR}"
    RESULT_VARIABLE UNZIP_RESULT
)
if(NOT UNZIP_RESULT EQUAL 0)
    message(FATAL_ERROR "Failed to unzip libproxy_core.zip")
endif()

add_library(${PLUGIN_NAME} SHARED
    "proxy_core_plugin.cpp"
    "proxy_core_plugin_c_api.cpp"
)

apply_standard_settings(${PLUGIN_NAME})
set_target_properties(${PLUGIN_NAME} PROPERTIES
    CXX_VISIBILITY_PRESET hidden
)
target_compile_definitions(${PLUGIN_NAME} PRIVATE FLUTTER_PLUGIN_IMPL)
target_include_directories(${PLUGIN_NAME} INTERFACE
    "${CMAKE_CURRENT_SOURCE_DIR}/include"
)
target_link_libraries(${PLUGIN_NAME} PRIVATE flutter flutter_wrapper_plugin)

# Bundle runtime libraries next to the host binary.
set(proxy_core_bundled_libraries
    "${PROXY_CORE_DLL}"
    PARENT_SCOPE
)

# If wintun.dll is shipped alongside, include it too.
set(WINTUN_DLL "${CMAKE_CURRENT_SOURCE_DIR}/wintun.dll")
if(EXISTS ${WINTUN_DLL})
    set(proxy_core_bundled_libraries
        "${PROXY_CORE_DLL}"
        "${WINTUN_DLL}"
        PARENT_SCOPE
    )
endif()
