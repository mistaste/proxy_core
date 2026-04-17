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

# Fetch wintun.dll (the userspace TUN driver from the WireGuard authors)
# once per configure if it's not already present. The SHA256 is the official
# wintun-0.14.1.zip hash from wintun.net and is enforced to guard against
# tampered mirrors; bump it in lock-step with the URL when upgrading.
set(WINTUN_DLL "${CMAKE_CURRENT_SOURCE_DIR}/wintun.dll")
set(WINTUN_URL "https://www.wintun.net/builds/wintun-0.14.1.zip")
set(WINTUN_SHA256 "07c256185d6ee3652e09fa55c0b673e2624b565e02c4b9091c79ca7d2f24ef51")

if(NOT EXISTS ${WINTUN_DLL})
    set(WINTUN_ZIP "${CMAKE_CURRENT_BINARY_DIR}/wintun.zip")
    set(WINTUN_STAGE "${CMAKE_CURRENT_BINARY_DIR}/wintun_stage")
    message(STATUS "Fetching wintun.dll from ${WINTUN_URL}")
    file(DOWNLOAD "${WINTUN_URL}" "${WINTUN_ZIP}"
        EXPECTED_HASH SHA256=${WINTUN_SHA256}
        SHOW_PROGRESS
        TLS_VERIFY ON
        STATUS WINTUN_STATUS)
    list(GET WINTUN_STATUS 0 WINTUN_STATUS_CODE)
    if(NOT WINTUN_STATUS_CODE EQUAL 0)
        message(FATAL_ERROR "Failed to download wintun.dll: ${WINTUN_STATUS}")
    endif()
    file(MAKE_DIRECTORY "${WINTUN_STAGE}")
    execute_process(
        COMMAND ${CMAKE_COMMAND} -E tar xf "${WINTUN_ZIP}"
        WORKING_DIRECTORY "${WINTUN_STAGE}")
    file(COPY "${WINTUN_STAGE}/wintun/bin/amd64/wintun.dll"
        DESTINATION "${CMAKE_CURRENT_SOURCE_DIR}")
    file(REMOVE_RECURSE "${WINTUN_STAGE}" "${WINTUN_ZIP}")
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

# Bundle runtime libraries next to the host binary. wintun.dll is always
# present at this point because the fetch block above would have failed
# the configure step otherwise.
set(proxy_core_bundled_libraries
    "${PROXY_CORE_DLL}"
    "${WINTUN_DLL}"
    PARENT_SCOPE
)
