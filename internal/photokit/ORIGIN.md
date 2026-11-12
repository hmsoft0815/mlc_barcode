# Origin

Copy of `photokit` from the OpticScript repository (`mlc_opticscript/photokit`,
module `gitlab.com/mlc0911-group/mlc_imgpipeline/photokit`), commit `7358814f`,
taken over unchanged on 2026-11-12 (ticket T-BAR-20261112-02).

It is copied rather than required because that repository is private and
mlc_barcode is public: a `replace` onto a sibling directory would break every
build that lacks it (GitHub CI, other clones). Changes belong upstream first;
then copy the files again and update the commit above.
