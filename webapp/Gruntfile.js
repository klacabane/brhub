module.exports = function(grunt) {
  grunt.initConfig({
    concat: {
      options: {
        separator: ';',
      },
      js: {
        src: ['js/helpers/*.js', 'js/models/*.js', 'js/modules/*.js', 'js/app.js'],
        dest: 'src/app.js',
      }
    },
    uglify: {
      js: {
        files: {
          'src/app.min.js': ['<%= concat.js.dest %>'],
        },
      },
    },
    clean: ['<%= concat.js.dest %>'],
    watch: {
      js: {
        files: 'js/**',
        tasks: ['default'],
      },
    },
  });

  grunt.loadNpmTasks('grunt-contrib-concat');
  grunt.loadNpmTasks('grunt-contrib-uglify');
  grunt.loadNpmTasks('grunt-contrib-watch');
  grunt.loadNpmTasks('grunt-contrib-clean');

  grunt.registerTask('default', ['concat', 'uglify', 'clean']);
};
