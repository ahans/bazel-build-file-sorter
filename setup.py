import setuptools

from setuptools.command.bdist_wheel import bdist_wheel


class patched_bdist_wheel(bdist_wheel):
    def finalize_options(self):
        bdist_wheel.finalize_options(self)
        # Mark us as not a pure python package
        self.root_is_pure = False

    def get_tag(self):
        _, _, plat = bdist_wheel.get_tag(self)
        # We don't contain any python source, nor any python extensions
        return 'py2.py3', 'none', plat


cmdclass = {'bdist_wheel': patched_bdist_wheel}

setuptools.setup(cmdclass=cmdclass)