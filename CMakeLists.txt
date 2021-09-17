cmake_minimum_required(VERSION 3.10)
project(rhosts VERSION 0.0.0 LANGUAGES C)
message("  Project Binary dir: '${PROJECT_SOURCE_DIR}'")
include_directories(${PROJECT_SOURCE_DIR})
add_executable(rhosts src/rhosts.c src/download.c)
target_link_libraries(rhosts curl)
target_include_directories(rhosts PUBLIC "${PROJECT_BINARY_DIR}")
configure_file(src/rhostsConfig.h.in config.h)
set(CMAKE_INSTALL_PREFIX "/usr/local")
install(TARGETS rhosts DESTINATION bin)
install(FILES "${PROJECT_SOURCE_DIR}/src/systemd/rhosts.service"
	DESTINATION rhosts/systemd)
install(FILES "${PROJECT_SOURCE_DIR}/src/systemd/rhosts.timer"
	DESTINATION rhosts/systemd)
install(FILES "${PROJECT_SOURCE_DIR}/src/systemd/rhosts.path"
	DESTINATION rhosts/systemd)
